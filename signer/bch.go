package signers

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	sigHashForkID = 0x40
)

// BchTxInputSignature defines input signature function for BTC
func BchTxInputSignature(tx *wire.MsgTx, idx int, subScript []byte,
	xpath []uint32, amount uint64, keyProvider KeyProvider) ([]byte, error) {
	hashType := txscript.SigHashAll
	hash, err := bip143SignatureHash(subScript, txscript.NewTxSigHashes(tx), hashType, tx, idx, amount, sigHashForkID)
	if err != nil {
		return nil, err
	}
	sig, err := keyProvider.SignDerived(hash, xpath)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	return append(sig, byte(hashType)|sigHashForkID), nil
}
