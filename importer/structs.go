package importer

import "github.com/wedancedalot/decimal"

type (
	NodeParams struct {
		Host     string
		Port     int
		User     string
		Password string
	}

	Operation struct {
		TxId         string
		TxOut        uint
		TxMemo       string
		ToAddress    string
		CurrencyCode string
		Amount       decimal.Decimal
	}
)
