package signers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"signer/conf"

	"github.com/apex/log"
	"github.com/stanche/crypto-interface/signer/script"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

var BtcNetParams = chaincfg.MainNetParams // TestNet3Params

type (

	// RawTxInputSignature defines signature function for BTC-class currencies
	RawTxInputSignature func(tx *wire.MsgTx, idx int, subScript []byte,
		xpath []uint32, amount uint64, keyProvider KeyProvider) ([]byte, error)

	// BtcSigner defines BTC-like signers
	BtcSigner struct {
		currency       string
		net            *chaincfg.Params
		keyProvider    KeyProvider
		inputSignature RawTxInputSignature
	}

	btcSignature struct {
		Ind int    `json:"i"`
		Val []byte `json:"v"`
	}

	// BtcPublicAttributes defines the return structure for btcSigner.Public()
	BtcPublicAttributes struct {
		XPub string
	}
)

// NewBtcSigner returns new instance of BtcSigner with the InputSignature function provided
func NewBtcSigner(currencyCode string, keyProvider KeyProvider, params *chaincfg.Params, inSign RawTxInputSignature) *BtcSigner {

	signer := BtcSigner{
		currency:       currencyCode,
		net:            params,
		keyProvider:    keyProvider,
		inputSignature: inSign,
	}
	return &signer
}

// BtcTxInputSignature defines input signature function for BTC
// RawTxInSignature returns the serialized ECDSA signature for the input idx of
// the given transaction, with hashType appended to it.
func BtcTxInputSignature(tx *wire.MsgTx, idx int, subScript []byte,
	xpath []uint32, amount uint64, keyProvider KeyProvider) ([]byte, error) {
	hashType := txscript.SigHashAll

	hash, err := txscript.CalcSignatureHash(subScript, hashType, tx, idx)
	if err != nil {
		return nil, err
	}
	sig, err := keyProvider.SignDerived(hash, xpath)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	return append(sig, byte(hashType)), nil
}

// CurrencyType implements Signer interface
func (signer *BtcSigner) CurrencyType() string {
	return signer.currency
}

// Sign implements Signer interface
func (signer *BtcSigner) Sign(txHex []byte, signParams []uint64) ([]string, error) {
	// parse the TxHex
	txData, err := hex.DecodeString(string(txHex))
	if err != nil {
		return nil, err
	}

	tx, err := btcutil.NewTxFromBytes(txData)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		err = fmt.Errorf("decoded tx is nil")
		return nil, err
	}

	xpubExp, err := signer.keyProvider.GetPublicKey()
	if err == nil && xpubExp == nil {
		err = fmt.Errorf("xpubExp is nil")
	}
	if err != nil {
		return nil, err
	}

	// extract TxIn count
	countTxIn := len(tx.MsgTx().TxIn)
	signatures := make([]string, countTxIn)
	//indexes := make([]uint64, countTxIn)

	flagParams := len(signParams) > 0
	if flagParams && len(signParams) != countTxIn {
		err = fmt.Errorf("inconsistnence signParams (%d) and countTxIn(%d)", len(signParams), countTxIn)
		return nil, err
	}
	// loop for TxIn's
	for indexTxIn := 0; indexTxIn < countTxIn; indexTxIn++ {

		// initialize dummy value
		//indexes[indexTxIn] = -1

		redeemScript, err := script.RedeemScriptFromTxin(tx.MsgTx().TxIn[indexTxIn])
		if err == nil && redeemScript == nil {
			err = fmt.Errorf("redeemScript is nil")
		}
		if err != nil {
			return nil, err
		}

		m, pubkeys, index, xpath, err := script.PubkeysIndexPathFromScript(redeemScript, xpubExp)
		if err == nil && pubkeys == nil {
			err = fmt.Errorf("pubkeys is nil")
		}
		if err != nil {
			return nil, err
		}
		// check the input belong to the key
		if index < 0 {
			continue
		}
		if xpath == nil {
			err = fmt.Errorf("xpath is nil")
			return nil, err
		}

		msScript, index, err := script.MultisigScriptFromPubkeys(m, pubkeys, index)
		if err == nil && msScript == nil {
			err = fmt.Errorf("multisig script is nil")
		}
		if err != nil {
			return nil, err
		}

		var amount uint64
		if flagParams {
			amount = signParams[indexTxIn]
		}
		sign, err := signer.inputSignature(tx.MsgTx(), indexTxIn, msScript, xpath, amount, signer.keyProvider)
		if err == nil && sign == nil {
			err = fmt.Errorf("sign is nil")
		}
		if err != nil {
			return nil, err
		}
		log.Debugf("sign[%d]: %s\n", indexTxIn, hex.EncodeToString(sign))

		btcSignature := btcSignature{
			Ind: index,
			Val: sign,
		}

		byteSignature, err := json.Marshal(btcSignature)
		if err != nil {
			return nil, err
		}
		signatures[indexTxIn] = base64.StdEncoding.EncodeToString(byteSignature)
		log.Debugf("sign[%d]: %s", indexTxIn, signatures[indexTxIn])
	}
	return signatures, nil
}

// Public returns Extended Public Key as string
func (signer *BtcSigner) Public() (interface{}, error) {
	net := &conf.BtcNetParams
	pk, err := signer.keyProvider.GetPublicKey()
	if err != nil {

		return "", err
	}
	if pk == nil {
		return "", fmt.Errorf("public key is nil")
	}

	pub := (*btcec.PublicKey)(pk)
	key := pub.SerializeCompressed()
	chainCode, err := signer.keyProvider.GetChainCode()
	if err != nil {

		return "", err
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}

	xpub := hdkeychain.NewExtendedKey(net.HDPublicKeyID[:], key, chainCode, parentFP, 0, 0, false)

	if xpub == nil {
		return "", fmt.Errorf("xpub is nil")
	}
	return BtcPublicAttributes{
		XPub: xpub.String(),
	}, nil
}
