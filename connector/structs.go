package connector

import (
	"math/big"

	"github.com/wedancedalot/decimal"
)

type (
	WalletParams struct {
		Currency       string
		Active         bool
		Type           string
		Node           NodeParams
		ChainConfig    string
		NodeTimeoutSec int
		FeeMax         int
		Core           NodeParams
		Debug          bool
	}
	// AddressBalance contains confirmed, unconfirmed and unmatured
	AddressBalance struct {
		Confirmed   decimal.Decimal
		Unconfirmed decimal.Decimal
		Unmatured   decimal.Decimal
	}

	// UtxoStruct defines return record from Utxos()
	UtxoStruct struct {
		TxHash  string
		Height  int
		TxPos   int
		Value   decimal.Decimal
		Address string
	}

	// WalletSignStruct - wallet parameters for Electrum
	WalletSignStruct struct {
		Signers uint8
		XPubs   []string
	}

	// UtxStruct defines Tx inputs for Electrum
	UtxStruct struct {
		TxHash string
		TxPos  int
		//Value  decimal.Decimal
		Index uint32
	}
	OutputParsed struct {
		Address string
		Value   *big.Int
		TxPos   uint
	}

	// OutStruct - defines outputs
	OutStruct struct {
		Address               string
		Amount                decimal.Decimal
		Currency              Currency
		Memo                  string
		SubtractFeeFromAmount bool
		IsChange              bool
	}

	// TxStatusStruct defines response from Connector.TxGet()
	// int64 used as answer could contain (-1, -1) in case of fork when the Tx was discarded
	TxStatusStruct struct {
		Height         int64
		Conf           uint64
		Fee            *big.Int
		IsIrreversible bool
	}

	Operation struct {
		TxId         string
		TxOut        uint
		TxTag        string
		ToAddress    string
		CurrencyCode string
		Amount       decimal.Decimal
	}
)

// NewTxStatusWithNonNeg creates a new TxStatusStruct gets number of confirmations. If the number of confirmations is negative returns zero.
func NewTxStatusWithNonNeg(h, c int64) TxStatusStruct {
	if c < 0 {
		c = 0
	}
	return TxStatusStruct{Height: h, Conf: uint64(c)}
}

type Connector struct {
	WalletId   uint64
	Currency   string
	WalletType string
}

func (c *Connector) CurrencyCode() string {
	return c.Currency
}

func (c *Connector) WalletID() uint64 {
	return c.WalletId
}

func (c *Connector) GetWalletType() string {
	return c.WalletType
}
