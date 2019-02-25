package hd

// GeneratorParameters defines parameters for Generator
type GeneratorParameters struct {
	SignersXpubs    []string
	SignersRequired uint8
	PathIndex       uint32
	Regtest         bool
}
