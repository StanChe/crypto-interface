package signers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stanche/crypto-interface/signer/script"
)

var masterKey = []byte("Bitcoin seed")

// Signer256k1 implements a key provider on sec256k1 curve.
type Signer256k1 struct {
	keyData   []byte
	ethKostil bool
}

// New creates a BTC-like key provider.
func New(secret []byte) Signer256k1 {
	return Signer256k1{keyData: secret}
}

// NewLegacyETH creates a key provider with a support for old ETH address generation.
func NewLegacyETH(secret []byte) Signer256k1 {
	return Signer256k1{keyData: secret, ethKostil: true}
}

func (s Signer256k1) getKey(path []uint32) (*btcec.PrivateKey, error) {
	if s.ethKostil && len(path) == 2 && path[0] == 0 && path[1] == 0 {
		key, err := ecdsa.GenerateKey(btcec.S256(), bytes.NewReader(s.keyData))
		if err != nil {
			return nil, err
		}
		return (*btcec.PrivateKey)(key), nil
	}
	// build the master key
	xkey, err := hdkeychain.NewMaster(s.keyData, &BtcNetParams)

	if err == nil && xkey == nil {
		err = fmt.Errorf("xkey is nil")
	}
	if err != nil {
		return nil, err
	}
	// clear the key
	defer xkey.Zero()

	xkeyIndex, err := script.ChildFromXkeyPath(xkey, path)
	if err == nil && xkeyIndex == nil {
		err = fmt.Errorf("child xkey is nil")
	}
	if err != nil {
		return nil, err
	}
	defer xkeyIndex.Zero()
	prvKey, err := xkeyIndex.ECPrivKey()
	if err == nil && prvKey == nil {
		err = fmt.Errorf("child key is nil")
	}
	if err != nil {
		return nil, err
	}
	return prvKey, err
}

func (s Signer256k1) signDerivedRS(hash []byte, path []uint32) (R, S *big.Int, err error) {

	prvKey, err := s.getKey(path)
	if err == nil && prvKey == nil {
		err = fmt.Errorf("child key is nil")
	}
	if err != nil {
		return nil, nil, err
	}
	// defer zeroKey(prvKey)
	signature, err := prvKey.Sign(hash)
	if err != nil {
		return nil, nil, err
	}
	return signature.R, signature.S, nil
}

// SignDerived signs the hash using the key on path in HD tree. If it's a signer with ethKostil - then [0,0] path is considered a special case.
func (s Signer256k1) SignDerived(hash []byte, path []uint32) ([]byte, error) {
	sr, ss, err := s.signDerivedRS(hash, path)
	if err != nil {
		return nil, err
	}
	signature := btcec.Signature{R: sr, S: ss}
	return signature.Serialize(), nil
}

// DerivedPubkey returns a ecdsa.PublicKey that relates to path on HD tree. 0,0] path is considered a special case for ethKostil.
func (s Signer256k1) DerivedPubkey(path []uint32) (*ecdsa.PublicKey, error) {
	prvKey, err := s.getKey(path)
	if err == nil && prvKey == nil {
		err = fmt.Errorf("child key is nil")
	}
	if err != nil {
		return nil, err
	}
	return &prvKey.PublicKey, nil

}

// GetPublicKey returns a master public key. It's generation depends on ethKostil.
func (s Signer256k1) GetPublicKey() (*ecdsa.PublicKey, error) {
	if s.ethKostil {
		key, err := ecdsa.GenerateKey(btcec.S256(), bytes.NewReader(s.keyData))
		if err != nil {
			return nil, err
		}
		// defer zeroKey(key)
		return &key.PublicKey, nil
	}

	// build the master key
	xkey, err := hdkeychain.NewMaster(s.keyData, &BtcNetParams)
	if err != nil {
		return nil, err
	}
	if xkey == nil {
		return nil, fmt.Errorf("xkey is nil")
	}

	// clear the key
	defer xkey.Zero()

	xpubExp, err := xkey.Neuter()
	if err != nil {
		return nil, err
	}
	if xpubExp == nil {
		return nil, fmt.Errorf("xpubExp is nil")
	}
	pk, err := xpubExp.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("ECPubKey from xpub failed: %s", err.Error())
	}
	if pk == nil {
		return nil, fmt.Errorf("ECPubKey returned nil")
	}
	return pk.ToECDSA(), err
}

// GetChainCode returns a deterministic chain code.
func (s Signer256k1) GetChainCode() ([]byte, error) {
	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(s.keyData)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	chainCode := lr[len(lr)/2:]
	return chainCode, nil
}
