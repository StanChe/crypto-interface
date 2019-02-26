package importer

import "github.com/wedancedalot/decimal"

type (
	Operation struct {
		TxId         string
		TxOut        uint
		TxMemo       string
		ToAddress    string
		CurrencyCode string
		Amount       decimal.Decimal
	}
)
