package btc_example

import (
	"fmt"
	"github.com/Nargott/goutils"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stanche/crypto-interface/connector"
	"github.com/wedancedalot/decimal"
	"math/big"
	"strings"
)

type (
	BtcBlockChainImporter struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}

	// outputParsed - describes return of TxParse
	outputParsed struct {
		Address string
		Value   *big.Int
		TxPos   uint
	}

	// processTxData is used for processTransaction func as param
	processTxData struct {
		block       *wire.MsgBlock
		txMsg       *wire.MsgTx
		blockNumber uint64
		currency    connector.Currency
		addresses   connector.AddressLister
	}

	// processTxData is used for processTransaction func as return
	processTxResponse struct {
		ops []connector.Operation
		err error
	}
)

var ErrBadCurrenciesCount = fmt.Errorf("bad currencies count provided: Bitcoin import was only supporting one currency BTC")

// NewBlockChainImporter creates new instance of importer.BlockChainImporter as BtcBlockChainImporter
func NewBlockChainImporter(node connector.NodeParams, chainParams chaincfg.Params, txBatchSize int) (connector.BlockChainImporter, error) {
	cl, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", node.GetHost(), node.GetPort()),
		User:         node.GetUser(),
		Pass:         node.GetPassword(),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Bitcoin node: %v", err.Error())
	}
	if cl == nil {
		return nil, fmt.Errorf("failed to connect to Bitcoin node: client is nil")
	}

	return BtcBlockChainImporter{
		client:      cl,
		chainParams: chainParams,
		txBatchSize: txBatchSize,
	}, nil
}

// getBlockByNumber returns btcd/wire MsgBlock as well
func (bci BtcBlockChainImporter) getBlockByNumber(number uint64) (block *wire.MsgBlock, err error) {
	if bci.client == nil {
		return nil, connector.ErrClientNil
	}

	blockHash, err := bci.client.GetBlockHash(int64(number))
	if blockHash == nil || err != nil {
		return nil, connector.ErrNotFound
	}

	block, err = bci.client.GetBlock(blockHash)
	if block == nil || err != nil {
		return nil, connector.ErrNotFound
	}

	return block, nil
}

// GetBlockHashesByNumber returns block hash and previous block gash as strings
func (bci BtcBlockChainImporter) GetBlockHashesByNumber(number uint64) (hash, prevHash string, err error) {
	block, err := bci.getBlockByNumber(number)
	if err != nil {
		return "", "", err
	}
	return block.Header.BlockHash().String(), block.Header.PrevBlock.String(), nil
}

// ProcessBlock do all importer logic and returns operations with given addresses list included in a given block
func (bci BtcBlockChainImporter) ProcessBlock(blockNumber uint64, currencies []connector.Currency, addresses connector.AddressLister) (operations []connector.Operation, err error) {
	//BTC importer only supports one currency at all
	if len(currencies) != 1 || strings.EqualFold(currencies[0].GetCode(), "BTC") {
		return operations, ErrBadCurrenciesCount
	}

	block, err := bci.getBlockByNumber(blockNumber)
	if err != nil {
		return operations, err
	}
	// Scan all transactions inside a block
	i := 0
	txCount := len(block.Transactions)
	var (
		lastResp processTxResponse
		errors   []error
	)
	for i < txCount {
		batchSize := goutils.Min(bci.txBatchSize, txCount-i)
		respCh := make(chan processTxResponse, batchSize)
		for txNumber := i; txNumber < i+batchSize; txNumber++ {
			txtoProcess := block.Transactions[txNumber]
			go func() {
				respCh <- bci.processTransaction(processTxData{
					block:       block,
					txMsg:       txtoProcess,
					blockNumber: blockNumber,
					currency:    currencies[0],
					addresses:   addresses,
				})
			}()
		}

		for txNumber := i; txNumber < i+batchSize; txNumber++ {
			lastResp = <-respCh
			if lastResp.err != nil {
				//collect errors
				errors = append(errors, fmt.Errorf("processTransaction [hash: %s] err: %s", block.Transactions[txNumber].TxHash().String(), lastResp.err.Error()))
			}
			operations = append(operations, lastResp.ops...)
		}

		i += batchSize
	}

	if len(errors) > 0 { //return all collected errors at once
		var sb strings.Builder
		for _, err := range errors {
			sb.WriteString(fmt.Sprintf("%#v \n", err))
		}
		return operations, fmt.Errorf(sb.String())
	}

	return operations, nil
}

// processTransaction returns operations on given addresses list, which included in transaction
func (bci BtcBlockChainImporter) processTransaction(d processTxData) processTxResponse {
	parsedOutputs, err := bci.parseOutputs(d.txMsg.TxOut)
	if err != nil {
		return processTxResponse{ops: nil, err: fmt.Errorf("btc processTransaction.ParseOutputs %s : %v", d.txMsg.TxHash(), err.Error())}
	}

	var operations []connector.Operation
	for _, output := range parsedOutputs {
		if d.addresses.HasAddress(output.Address, "") {
			operations = append(operations, connector.Operation{
				TxId:      d.txMsg.TxHash().String(),
				TxOut:     output.TxPos,
				ToAddress: output.Address,
				Amount:    decimal.NewFromBigInt(output.Value, -int32(d.currency.GetPrecision())),
			})
		}
	}

	return processTxResponse{ops: operations, err: err}
}

// parseOutputs parses all BTC outputs and returns in convenient format outputParsed
func (bci BtcBlockChainImporter) parseOutputs(txOuts []*wire.TxOut) ([]outputParsed, error) {
	var outputs []outputParsed
	for i, txOut := range txOuts {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, &bci.chainParams)
		if err != nil {
			//todo find out what to do in case if we cannot parse output address
			continue
		}

		for _, address := range addresses {
			outputs = append(outputs, outputParsed{
				Address: address.EncodeAddress(),
				Value:   big.NewInt(txOut.Value),
				TxPos:   uint(i),
			})
		}
	}

	return outputs, nil
}
