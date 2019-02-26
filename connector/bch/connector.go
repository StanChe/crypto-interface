package bch

import (
	"fmt"
	"math/big"

	bchchaincfg "github.com/bchsuite/bchd/chaincfg"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	"github.com/Messer4/bchaddr"
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/stanche/crypto-interface/connector"
	"github.com/stanche/crypto-interface/connector/btc"
)

type (
	bchChainConnector struct {
		connector.Connector
		ibtc    btc.IBtcChainConnector
		chain   *chaincfg.Params
		regtest bool
	}
)

// NewChainConnector returns the IBtcChainConnector interface
// to use btc importer as ltc one
func NewChainConnector(walletID uint64, cfg *connector.WalletParams) (btc.IBtcChainConnector, error) {
	iBtcConnector, err := btc.NewBtcChainConnector(walletID, cfg)
	if err != nil {
		return nil, err
	}

	connector := &bchChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   cfg.Currency,
			WalletType: cfg.Type,
		},
		ibtc:    iBtcConnector,
		chain:   &bchchaincfg.TestNet3Params,
		regtest: cfg.ChainConfig == "regtest",
	}
	connector.ibtc.DecoderSet(connector.DecodeAddress)

	return connector, nil
}

func (c *bchChainConnector) BalanceGet(_ connector.Currency, address ...string) (balance connector.AddressBalance, err error) {
	// log.Warnf("bchChainConnector does not support BalanceGet method")
	return balance, fmt.Errorf("unsupported method: BalanceGet")
}

func (c *bchChainConnector) ValidateAddress(address string) (bool, error) {

	_, err := bchaddr.ToLegacyAddress(address)
	return err == nil, nil
}

func (c *bchChainConnector) DecodeAddress(addr string) (btcutil.Address, error) {

	legacyAddress, err := bchaddr.ToLegacyAddress(addr)
	if err != nil {
		return nil, err
	}

	bchAddr := btc.CoinAddress{
		Addr: legacyAddress,
	}
	return &bchAddr, err
}

func (c *bchChainConnector) DecoderSet(decoder btc.AddressDecoder) {
	//
}

func (c *bchChainConnector) CreateRawTransaction(inputs []btcjson.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {

	return c.ibtc.CreateRawTransaction(inputs, amounts)
}

func (c *bchChainConnector) ParseOutputs(txOuts []*wire.TxOut) ([]*connector.OutputParsed, error) {
	var outputs []*connector.OutputParsed
	for index := range txOuts {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txOuts[index].PkScript, c.chain)
		if err != nil {
			// log.Errorf("ExtractTxOutAddresses %s", err.Error())
			return nil, err
		}

		for _, address := range addresses {
			addr, err := bchaddr.ToCashAddress(address.EncodeAddress(), c.regtest)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, &connector.OutputParsed{
				Address: addr,
				Value:   big.NewInt(txOuts[index].Value),
				TxPos:   uint(index),
			})
		}
	}
	return outputs, nil
}

func (c *bchChainConnector) GetBlockByNumber(number uint64) (*wire.MsgBlock, error) {
	return c.ibtc.GetBlockByNumber(number)
}

func (c *bchChainConnector) GetTransactionByHash(hash chainhash.Hash) (*btcjson.TxRawResult, bool, error) {
	return c.ibtc.GetTransactionByHash(hash)
}

// TxStatus returns transaction status by TxId(hash)
func (c *bchChainConnector) TxStatus(txID string, blockNo uint64) (*connector.TxStatusStruct, error) {
	return c.ibtc.TxStatus(txID, blockNo)
}

func (c *bchChainConnector) TxBuild(walletData *connector.WalletSignStruct,
	utxosIn interface{}, output []connector.OutStruct) (string, error) {

	return c.ibtc.TxBuild(walletData, utxosIn, output)
}

func (c *bchChainConnector) TxBroadcast(txHex string) (string, error) {
	return c.ibtc.TxBroadcast(txHex)
}

// TxRebuild - combine parsed hex Tx with the signatures
func (c *bchChainConnector) TxRebuild(txHex string, signatures connector.TxSignatures) (string, error) {
	return btc.TxRebuildBtc(txHex, signatures)
}
