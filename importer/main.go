package importer

import (
	"github.com/stanche/crypto-interface/connector"
)

type (
	NodeParams interface {
		GetHost() string
		GetPort() int
		GetUser() string
		GetPassword() string
	}

	BlockChainImporter interface {
		GetBlockHashesByNumber(number uint64) (hash, prevHash string, err error)
		ProcessBlock(blockNumber uint64, currencies []connector.Currency, addresses []connector.Address) (operations []Operation, err error)
	}
)
