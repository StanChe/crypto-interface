package signers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

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
		key *btcec.PrivateKey, amount uint64) ([]byte, error)

	// BtcSigner defines BTC-like signers
	BtcSigner struct {
		currency       string
		net            *chaincfg.Params
		keyData        []byte
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
func NewBtcSigner(currencyCode string, secret []byte, params *chaincfg.Params, inSign RawTxInputSignature) *BtcSigner {

	signer := BtcSigner{
		currency:       currencyCode,
		net:            params,
		keyData:        secret,
		inputSignature: inSign,
	}
	return &signer
}

// BtcTxInputSignature defines input signature function for BTC
func BtcTxInputSignature(tx *wire.MsgTx, idx int, subScript []byte,
	prvKey *btcec.PrivateKey, amount uint64) ([]byte, error) {

	return txscript.RawTxInSignature(tx, idx, subScript, txscript.SigHashAll, prvKey)
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

	// build the master key
	xkey, err := hdkeychain.NewMaster(signer.keyData, &BtcNetParams)
	if err == nil && xkey == nil {
		err = fmt.Errorf("xkey is nil")
	}
	if err != nil {
		return nil, err
	}
	// clear the key
	defer xkey.Zero()

	xpubExp, err := xkey.Neuter()
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

		xkeyIndex := xkey
		xkeyIndex, err = script.ChildFromXkeyPath(xkeyIndex, xpath)
		if err == nil && xkeyIndex == nil {
			err = fmt.Errorf("child xkey is nil")
		}
		if err != nil {
			return nil, err
		}
		prvKey, err := xkeyIndex.ECPrivKey()
		if err == nil && prvKey == nil {
			err = fmt.Errorf("child key is nil")
		}
		if err != nil {
			return nil, err
		}
		var amount uint64
		if flagParams {
			amount = signParams[indexTxIn]
		}
		sign, err := signer.inputSignature(tx.MsgTx(), indexTxIn, msScript, prvKey, amount)
		if err == nil && sign == nil {
			err = fmt.Errorf("sign is nil")
		}
		if err != nil {
			return nil, err
		}

		btcSignature := btcSignature{
			Ind: index,
			Val: sign,
		}

		byteSignature, err := json.Marshal(btcSignature)
		if err != nil {
			return nil, err
		}
		signatures[indexTxIn] = base64.StdEncoding.EncodeToString(byteSignature)
	}
	return signatures, nil
}

// Public returns Extended Public Key as string
func (signer *BtcSigner) Public() (interface{}, error) {

	net := &BtcNetParams

	// build the master key
	xkey, err := hdkeychain.NewMaster(signer.keyData, net)
	if err == nil && xkey == nil {
		err = fmt.Errorf("xkey is nil")
	}
	if err != nil {
		return "", err
	}
	// clear the key
	defer xkey.Zero()

	xpub, err := xkey.Neuter()
	if err == nil && xpub == nil {
		err = fmt.Errorf("xpub is nil")
	}
	if err != nil {
		return "", err
	}
	return BtcPublicAttributes{
		XPub: xpub.String(),
	}, nil
}
