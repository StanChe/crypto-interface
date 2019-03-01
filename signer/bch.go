package signers

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	sigHashForkID = 0x40
)

// BchTxInputSignature defines input signature function for BTC
func BchTxInputSignature(tx *wire.MsgTx, idx int, subScript []byte,
	prvKey *btcec.PrivateKey, amount uint64) ([]byte, error) {

	return bip143TxInSignature(tx, idx, subScript, txscript.SigHashAll, prvKey, amount, sigHashForkID)
}
