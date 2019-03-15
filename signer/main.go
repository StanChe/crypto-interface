package signers

import "crypto/ecdsa"

type TxSigner interface {
	CurrencyType() string
	Public() (interface{}, error)
	Sign(txData []byte, parameters []uint64) ([]string, error)
}

// KeyProvider provides a signing interface compatible with ledger.
type KeyProvider interface {
	// SignDerived signs the hash using the path in HD tree. The resulting signature is in a strict DER format:
	// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
	// Note that the serialized bytes returned do not include the appended hash type
	// used in Bitcoin signature scripts.
	SignDerived(hash []byte, path []uint32) ([]byte, error)

	// Sign(hash []byte) ([]byte, error)

	GetPublicKey() (*ecdsa.PublicKey, error)
	GetChainCode() ([]byte, error)
	DerivedPubkey(path []uint32) (*ecdsa.PublicKey, error)
}
