package importer

import (
	"github.com/stanche/crypto-interface/connectors"
)

type (
	BlockChainImporter interface {
		GetBlockHashesByNumber(number uint64) (hash, prevHash string, err error)
		ProcessBlock(blockNumber uint64, currencies []connectors.Currency, addresses []connectors.Address) (operations []Operation, err error)
	}
)
