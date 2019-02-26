package btc_example

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/fanliao/go-promise"
	"github.com/stanche/crypto-interface/connectors"
	"github.com/stanche/crypto-interface/importer"
	"github.com/wedancedalot/decimal"
	"math/big"
	"reflect"
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
		currency    connectors.Currency
		addresses   []connectors.Address
	}
)

var ErrBadCurrenciesCount = fmt.Errorf("bad currencies count provided: Bitcoin import was only supporting one currency BTC")

// NewBlockChainImporter creates new instance of importer.BlockChainImporter as BtcBlockChainImporter
func NewBlockChainImporter(node importer.NodeParams, chainParams chaincfg.Params, txBatchSize int) (importer.BlockChainImporter, error) {
	cl, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", node.Host, node.Port),
		User:         node.User,
		Pass:         node.Password,
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
		return nil, connectors.ErrClientNil
	}

	blockHash, err := bci.client.GetBlockHash(int64(number))
	if blockHash == nil || err != nil {
		return nil, connectors.ErrNotFound
	}

	block, err = bci.client.GetBlock(blockHash)
	if block == nil || err != nil {
		return nil, connectors.ErrNotFound
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
func (bci BtcBlockChainImporter) ProcessBlock(blockNumber uint64, currencies []connectors.Currency, addresses []connectors.Address) (operations []importer.Operation, err error) {
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
	for i < txCount {
		batchSize := bci.min(bci.txBatchSize, txCount-i)
		tasks := make([]interface{}, batchSize)
		taskNumber := 0
		for txNumber := i; txNumber < i+batchSize; txNumber++ {
			txtoProcess := block.Transactions[txNumber]
			tasks[taskNumber] = func() (r interface{}, err error) {
				ops, err := bci.processTransaction(processTxData{
					block:       block,
					txMsg:       txtoProcess,
					blockNumber: blockNumber,
					currency:    currencies[0],
					addresses:   addresses,
				})
				return ops, err
			}
			taskNumber++
		}

		f := promise.WhenAll(tasks...)
		result, err := f.Get()
		if err != nil {
			e, ok := err.(*promise.AggregateError)
			if !ok {
				return operations, fmt.Errorf("unexpected type of error: expected *promise.AggregateError, but received: %s", reflect.TypeOf(err))
			}
			return operations, e.InnerErrs[0]
		}

		ops, ok := result.([]importer.Operation)
		if !ok {
			return operations, fmt.Errorf("unexpected type of result: expected []importers.Operation, but received: %s", reflect.TypeOf(result))
		}

		operations = append(operations, ops...)
		i = i + batchSize
	}

	return operations, nil
}

// min returns a minimal value of a and b
func (bci BtcBlockChainImporter) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (bci BtcBlockChainImporter) isAddressInList(target string, addresses []connectors.Address) bool {
	for _, a := range addresses {
		if strings.EqualFold(target, a.GetAddress()) {
			return true
		}
	}

	return false
}

// processTransaction returns operations on given addresses list, which included in transaction
func (bci BtcBlockChainImporter) processTransaction(d processTxData) (operations []importer.Operation, err error) {
	parsedOutputs, err := bci.parseOutputs(d.txMsg.TxOut)
	if err != nil {
		return operations, fmt.Errorf("btc processTransaction.ParseOutputs %s : %v", d.txMsg.TxHash(), err.Error())
	}
	for _, output := range parsedOutputs {
		if bci.isAddressInList(output.Address, d.addresses) {
			operations = append(operations, importer.Operation{
				TxId:      d.txMsg.TxHash().String(),
				TxOut:     output.TxPos,
				ToAddress: output.Address,
				Amount:    decimal.NewFromBigInt(output.Value, -int32(d.currency.GetPrecision())),
			})
		}
	}

	return operations, nil
}

// parseOutputs parses all BTC outputs and returns in convenient format outputParsed
func (bci BtcBlockChainImporter) parseOutputs(txOuts []*wire.TxOut) ([]outputParsed, error) {
	var outputs []outputParsed
	for i, txOut := range txOuts {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, &bci.chainParams)
		if err != nil {
			//TODO find out what to do in case if we cannot parse output address
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
