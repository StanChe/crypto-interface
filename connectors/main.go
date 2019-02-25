package connectors

import (
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
	Currency              Currency
	Memo                  string
	SubtractFeeFromAmount bool
	IsChange              bool
}

// TxStatusStruct defines response from Connector.TxGet()
// int64 used as answer could contain (-1, -1) in case of fork when the Tx was discarded
type TxStatusStruct struct {
	Height         int64
	Conf           uint64
	Fee            *big.Int
	IsIrreversible bool
}

// NewTxStatusWithNonNeg creates a new TxStatusStruct gets number of confirmations. If the number of confirmations is negative returns zero.
func NewTxStatusWithNonNeg(h, c int64) TxStatusStruct {
	if c < 0 {
		c = 0
	}
	return TxStatusStruct{Height: h, Conf: uint64(c)}
}

// TxInSignatures - signatures array for a tx input
type TxInSignatures []string

// TxSignatures - signatures set for the tx
type TxSignatures []TxInSignatures

type (
	// Currency provides info about the currency.
	Currency interface {
		// GetCode is the code of the currency (i.e. BTC/ETH/USDT). It's usually capitalized.
		GetCode() string
		// GetPrecision gets the maximum number of decimal points of the currency.
		GetPrecision() uint8
		// GetTokenAddress gets the address of the contract of the currency (if applicable). If not applicable - returns an emppty string.
		GetTokenAddress() string
		// GetTokenCode gets an integer code of the currency (if applicable). If not applicable - returns 0.
		GetTokenCode() int64
	}

	// BalanceProvider is an interface for getting sum of the balances on specified addresses.
	BalanceProvider interface {
		BalanceGet(currency Currency, address ...string) (AddressBalance, error)
	}
	// AddressValidator is an interface for verifying whether the address is valid for the blockchain.
	// If it returns true - we are safe to send coins to that address.
	AddressValidator interface {
		ValidateAddress(address string) (bool, error)
	}

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
		BalanceProvider
		AddressValidator
		CurrencyCode() string
		WalletID() uint64
		GetWalletType() string
	}

	Connector struct {
		WalletId   uint64
		Currency   string
		WalletType string
	}
)

// TxPermanentFailure indicates a permanent failure for a tx - their is no need to try to broadcast the tx that returns such error.
var TxPermanentFailure = fmt.Errorf("transaction failed permanently")

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
