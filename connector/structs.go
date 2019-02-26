package connector

import (
	"math/big"

	"github.com/wedancedalot/decimal"
)

type (
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
)

// NewTxStatusWithNonNeg creates a new TxStatusStruct gets number of confirmations. If the number of confirmations is negative returns zero.
func NewTxStatusWithNonNeg(h, c int64) TxStatusStruct {
	if c < 0 {
		c = 0
	}
	return TxStatusStruct{Height: h, Conf: uint64(c)}
}
