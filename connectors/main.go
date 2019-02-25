package connectors

import (
	"cryptagio-walle/dao/models"
	"fmt"
	"math/big"

	"github.com/wedancedalot/decimal"
)

// AddressBalance contains confirmed, unconfirmed and unmatured
type AddressBalance struct {
	Confirmed   decimal.Decimal
	Unconfirmed decimal.Decimal
	Unmatured   decimal.Decimal
}

// Add creates a new AddressBalance instance with the sum of the balances
func (a AddressBalance) Add(b AddressBalance) AddressBalance {
	return AddressBalance{
		Confirmed:   a.Confirmed.Add(b.Confirmed),
		Unconfirmed: a.Unconfirmed.Add(b.Unconfirmed),
		Unmatured:   a.Unmatured.Add(b.Unmatured),
	}
}

// UtxoStruct defines return record from Utxos()
type UtxoStruct struct {
	TxHash  string
	Height  int
	TxPos   int
	Value   decimal.Decimal
	Address string
}

// WalletSignStruct - wallet parameters for Electrum
type WalletSignStruct struct {
	Signers uint8
	XPubs   []string
}

// UtxStruct defines Tx inputs for Electrum
type UtxStruct struct {
	TxHash string
	TxPos  int
	//Value  decimal.Decimal
	Index uint32
}

// OutStruct - defines outputs
type OutStruct struct {
	Address               string
	Amount                decimal.Decimal
	Currency              *models.Currency
	Memo                  string
	SubtractFeeFromAmount bool
	IsChange              bool
}

// InputParsed - describes return of TxParse
type InputParsed struct {
	PrevOutHash string
	PrevOutPos  int
	Sequence    uint
	XPubkeys    []*string
	//Pubkeys []*string
	Type         string
	Signers      int
	ScriptSig    string
	RedeemScript string
	//Signatures []*string
}

// OutputParsed - describes return of TxParse
type OutputParsed struct {
	Address string
	Value   *big.Int
	TxPos   uint
}

// TxStatusStruct defines response from Connector.TxGet()
// int64 used as answer could contain (-1, -1) in case of fork when the Tx was discarded
type TxStatusStruct struct {
	Height         int64
	Conf           uint64
	Fee            *big.Int
	IsIrreversible bool
}

// TxInSignatures - signatures array for a tx input
type TxInSignatures []string

// TxSignatures - signatures set for the tx
type TxSignatures []TxInSignatures

type (
	TxBuilder interface {
		TxBuild(walletData *WalletSignStruct, utxos interface{}, output []*OutStruct) (string, error)
		TxRebuild(txHex string, signatures TxSignatures) (string, error)
	}
	TxGetter interface {
		TxStatus(txID string, blockNo uint64) (*TxStatusStruct, error)
	}
	TxSender interface {
		TxBroadcast(txHex string) (txHash string, err error)
	}

	// IConnector defines wallet (node) interface
	IConnector interface {
		TxGetter
		TxBuilder
		TxSender
		CurrencyCode() string
		WalletID() uint64
		GetWalletType() string
		BalanceGet(currency models.Currency, address ...string) (*AddressBalance, error)
		ValidateAddress(address string) (bool, error)
	}

	Connector struct {
		WalletId   uint64
		Currency   string
		WalletType string
	}
)

var TxPermanentFailure = fmt.Errorf("transaction filed permanently")

func (c *Connector) CurrencyCode() string {
	return c.Currency
}

func (c *Connector) WalletID() uint64 {
	return c.WalletId
}

func (c *Connector) GetWalletType() string {
	return c.WalletType
}

var ErrNotFound = fmt.Errorf("not found")
