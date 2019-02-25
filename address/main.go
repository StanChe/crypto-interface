package address

import (
	"github.com/stanche/crypto-interface/address/hd"
)

type (
	// Generator defines generic interface for HD-adresses generation for deposits.
	Generator interface {
		AddressGenerate(params hd.GeneratorParameters) (addr string, err error)
	}
	// TagGenerator defines an interface for tag generation for supporting currencies.
	TagGenerator interface {
		TagGenerate() (tag string, err error)
	}

	// PublicKeyGenerator defines an interface to get a public key based on HD pattern.
	PublicKeyGenerator interface {
		PublicKeyGenerate(params hd.GeneratorParameters) (pubKey string, err error)
	}
)
