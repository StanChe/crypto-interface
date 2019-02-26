package btc

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stanche/crypto-interface/connector"
	"github.com/wedancedalot/decimal"
)

type (
	// OutputParsed - describes return of TxParse
	OutputParsed struct {
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
		addresses   []connector.Address
	}

	// processTxData is used for processTransaction func as return
	processTxResponse struct {
		ops []connector.Operation
		err error
	}
)

var ErrBadCurrenciesCount = fmt.Errorf("bad currencies count provided: Bitcoin import was only supporting one currency BTC")

// GetBlockHashesByNumber returns block hash and previous block gash as strings
func (bcc BtcChainConnector) GetBlockHashesByNumber(number uint64) (hash, prevHash string, err error) {
	block, err := bcc.GetBlockByNumber(number)
	if err != nil {
		return "", "", err
	}
	return block.Header.BlockHash().String(), block.Header.PrevBlock.String(), nil
}

// ProcessBlock do all importer logic and returns operations with given addresses list included in a given block
func (bcc BtcChainConnector) ProcessBlock(blockNumber uint64, currencies []connector.Currency, addresses []connector.Address) (operations []connector.Operation, err error) {
	//BTC importer only supports one currency at all
	if len(currencies) != 1 || strings.EqualFold(currencies[0].GetCode(), "BTC") {
		return operations, ErrBadCurrenciesCount
	}

	block, err := bcc.GetBlockByNumber(blockNumber)
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
		batchSize := bcc.min(bcc.txBatchSize, txCount-i)
		respCh := make(chan processTxResponse, batchSize)
		for txNumber := i; txNumber < i+batchSize; txNumber++ {
			txtoProcess := block.Transactions[txNumber]
			go func() {
				respCh <- bcc.processTransaction(processTxData{
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

// min returns a minimal value of a and b
func (BtcChainConnector) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isAddressInList returns true if target address in given addresses
func (BtcChainConnector) isAddressInList(target string, addresses []connector.Address) bool {
	for _, a := range addresses {
		if strings.EqualFold(target, a.GetAddress()) {
			return true
		}
	}

	return false
}

// processTransaction returns operations on given addresses list, which included in transaction
func (bcc BtcChainConnector) processTransaction(d processTxData) processTxResponse {
	parsedOutputs, err := bcc.parseOutputs(d.txMsg.TxOut)
	if err != nil {
		return processTxResponse{ops: nil, err: fmt.Errorf("btc processTransaction.ParseOutputs %s : %v", d.txMsg.TxHash(), err.Error())}
	}

	var operations []connector.Operation
	for _, output := range parsedOutputs {
		if bcc.isAddressInList(output.Address, d.addresses) {
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

// parseOutputs parses all BTC outputs and returns in convenient format OutputParsed
func (bcc BtcChainConnector) parseOutputs(txOuts []*wire.TxOut) ([]OutputParsed, error) {
	var outputs []OutputParsed
	for i, txOut := range txOuts {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, bcc.chain)
		if err != nil {
			//TODO find out what to do in case if we cannot parse output address
			continue
		}

		for _, address := range addresses {
			outputs = append(outputs, OutputParsed{
				Address: address.EncodeAddress(),
				Value:   big.NewInt(txOut.Value),
				TxPos:   uint(i),
			})
		}
	}

	return outputs, nil
}
