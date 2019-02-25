package btc

import (
	"crypto/ecdsa"
	"fmt"
	"sort"

	btcchaincfg "github.com/btcsuite/btcd/chaincfg"
	btscript "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	btckeychain "github.com/btcsuite/btcutil/hdkeychain"

	"github.com/stanche/crypto-interface/address/hd"
)

const MaxSigners = 15

type (
	// Generator is a struct able to create a trx address.
	Generator struct{}
)

// New creates a new trx generator instance
func New() Generator {
	return Generator{}
}

// AddressGenerate - main function for wallet service address generation
func (g Generator) AddressGenerate(params hd.GeneratorParameters) (address string, err error) {
	netParams := btcchaincfg.TestNet3Params
	return g.AddressGenerateForNet(params, netParams)
}

func (Generator) AddressGenerateForNet(params hd.GeneratorParameters, netParams btcchaincfg.Params) (address string, err error) {
	//func Generate(hdPath string, xPubs []string, signersRequired uint8) (address string, tag string, err error) {

	signersTotal := len(params.SignersXpubs)

	if signersTotal < 1 || signersTotal > MaxSigners {
		return "", fmt.Errorf("invalid signers quantity")
	}
	signersRequired := int(params.SignersRequired)
	if signersRequired < 1 || signersRequired > signersTotal {
		return "", fmt.Errorf("Invalid signersRequired: %d", signersRequired)
	}

	for i := 0; i < signersTotal; i++ {
		if params.SignersXpubs[i] == "" {
			return "", fmt.Errorf("invalid xpubs")
		}
	}

	//var xPub string
	//
	var btcExtKey *btckeychain.ExtendedKey

	hdPath := []uint32{
		0, params.PathIndex,
	}
	xPub, err := hd.XPubByPath(params.SignersXpubs[0], hdPath)
	if err != nil {

		return "", fmt.Errorf("xPubByPath error: %s", err.Error())
	}
	//var public *ecdsa.PublicKey
	_, btcExtKey, err = hd.PubKeyFromXpub(xPub)
	if err != nil {
		return "", err
	}

	if signersTotal == 1 {
		var btcAddressObj *btcutil.AddressPubKeyHash

		btcAddressObj, err = btcExtKey.Address(&netParams)
		if err != nil {
			return "", err
		}
		address = btcAddressObj.String()
	} else {
		var pubKeysList hd.BtcPubkeyList
		var publicECDSA *ecdsa.PublicKey
		var addressPubKey *btcutil.AddressPubKey
		var multisigScript []byte
		var addrScriptHash *btcutil.AddressScriptHash
		for i := 0; i < signersTotal; i++ {
			publicECDSA, err = hd.XPublicByHdPath(params.SignersXpubs[i], hdPath)
			if err != nil {
				return "", fmt.Errorf("GenerateKeyPairByHDKey error: %s", err.Error())
			}
			addressPubKey, err = btcutil.NewAddressPubKey(
				hd.BtcecPublicKeyFromECDSA(publicECDSA).SerializeCompressed(),
				&netParams)
			if err != nil {
				return "", fmt.Errorf("NewAddressPubKey error: %s", err.Error())
			}
			pubKeysList = append(pubKeysList, addressPubKey)
		}

		sort.Sort(pubKeysList)

		multisigScript, err = btscript.MultiSigScript(pubKeysList, int(signersRequired))
		if err != nil {

			return "", fmt.Errorf("MultiSigScript error: %s", err.Error())
		}

		addrScriptHash, err = btcutil.NewAddressScriptHash(multisigScript, &netParams)
		if err != nil {

			return "", fmt.Errorf("NewAddressScriptHash error: %s", err.Error())
		}
		address = addrScriptHash.String()
	}
	return address, nil
}
