package hd

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcutil"
	btckeychain "github.com/btcsuite/btcutil/hdkeychain"
	"github.com/xbis/godash/btcec"
)

func XPubByPath(xKeyMaster string, hdPath []uint32) (xPub string, err error) {

	xPub = xKeyMaster
	for i := range hdPath {

		xPub, err = childXpub(xPub, hdPath[i])
		if err != nil {
			return "", fmt.Errorf("childXpub error: %s", err.Error())

		}
	}

	return xPub, nil
}

func childXpub(xPubKey string, index uint32) (xPub string, err error) {
	extKey, err := btckeychain.NewKeyFromString(xPubKey)
	if err != nil {
		return "", fmt.Errorf("hdkeychain.NewKeyFromString (xkey = %s ) error: %s", xPubKey, err.Error())

	}

	childKey, err := extKey.Child(index)
	if err != nil {
		return "", fmt.Errorf("Child generation error: %s", err.Error())

	}

	xPub = childKey.String()
	return xPub, nil
}

func PubKeyFromXpub(xKey string) (publicKey *ecdsa.PublicKey, extKey *btckeychain.ExtendedKey, err error) {
	extKey, err = btckeychain.NewKeyFromString(xKey)
	if err != nil {
		return
	}

	childAddressObj, err := extKey.ECPubKey()
	if err != nil {
		return
	}
	publicKey = childAddressObj.ToECDSA()
	return
}

func XPublicByHdPath(xPublicKey string, hdPath []uint32) (publicKey *ecdsa.PublicKey, err error) {
	xPub, err := XPubByPath(xPublicKey, hdPath)
	if err != nil {
		return nil, fmt.Errorf("xPubByPath error: %s", err.Error())

	}
	publicKey, _, err = PubKeyFromXpub(xPub)
	return
}

func BtcecPublicKeyFromECDSA(ECSDAPublicKey *ecdsa.PublicKey) *btcec.PublicKey {
	pub := &btcec.PublicKey{
		Curve: ECSDAPublicKey.Curve,
		X:     ECSDAPublicKey.X,
		Y:     ECSDAPublicKey.Y,
	}
	return pub
}

type (
	BtcPubkeyList []*btcutil.AddressPubKey
)

func (a BtcPubkeyList) Len() int {
	return len(a)
}
func (a BtcPubkeyList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a BtcPubkeyList) Less(i, j int) bool {
	return bytes.Compare(a[i].ScriptAddress(), a[j].ScriptAddress()) == -1
}
