package signers

type TxSigner interface {
	CurrencyType() string
	Public() (interface{}, error)
	Sign(txData []byte, parameters []uint64) ([]string, error)
}
