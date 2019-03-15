package script

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
)

const (
	sizeEncodedKey = 4 + 1 + 4 + 4 + 32 + 33
	sizeCheckSum   = 4

	// at least 3 items shall be in the TxIn sigScript(multisig):
	// 0 (due to OP_CHECKMULTISIG), FF (NO_SIGNATURE), [..] (RedeemScript)
	sizeTxInMultisigScriptMin = 3

	sizeMultisigScriptMin     = 1 + 2 + 2
	sizeMultisigScriptMax     = 1 + 15 + 2
	sizeMultisigScriptMNCheck = 3
	sizeElectrumXpubHD        = 1 + sizeEncodedKey + 4 + 4
	markerXpub                = 0xff
)

// parseScript preparses the script in bytes into a list of parsedOpcodes while
// applying a number of sanity checks.  When there are parse errors, it returns
// the list of parsed opcodes up to the point of failure along with the error.
func parseScript(script []byte) ([]ParsedOpcode, error) {
	retScript := make([]ParsedOpcode, 0, len(script))
	for i := 0; i < len(script); {
		instr := script[i]
		op := &opcodes[instr]
		pop := ParsedOpcode{Opcode: op}

		// Parse data out of instruction.
		switch {
		// No additional data.  Note that some of the opcodes, notably
		// OP_1NEGATE, OP_0, and OP_[1-16] represent the data
		// themselves.
		case op.Length == 1:
			i++

		// Data pushes of specific lengths -- OP_DATA_[1-75].
		case op.Length > 1:
			if len(script[i:]) < op.Length {
				err := fmt.Errorf("opcode %02x requires %d "+
					"bytes, but script only has %d remaining",
					op.Value, op.Length, len(script[i:]))
				return retScript, err
			}

			// Slice out the data.
			pop.Data = script[i+1 : i+op.Length]
			i += op.Length

		// Data pushes with parsed lengths -- OP_PUSHDATAP{1,2,4}.
		case op.Length < 0:
			var l uint
			off := i + 1

			if len(script[off:]) < -op.Length {
				err := fmt.Errorf("opcode %02x requires %d "+
					"bytes, but script only has %d remaining",
					op.Value, -op.Length, len(script[off:]))
				return retScript, err
			}

			// Next -length bytes are little endian length of data.
			switch op.Length {
			case -1:
				l = uint(script[off])
			case -2:
				l = ((uint(script[off+1]) << 8) |
					uint(script[off]))
			case -4:
				l = ((uint(script[off+3]) << 24) |
					(uint(script[off+2]) << 16) |
					(uint(script[off+1]) << 8) |
					uint(script[off]))
			default:
				err := fmt.Errorf("invalid opcode length %d",
					op.Length)
				return retScript, err
			}

			// Move offset to beginning of the data.
			off += -op.Length

			// Disallow entries that do not fit script or were
			// sign extended.
			if int(l) > len(script[off:]) || int(l) < 0 {
				err := fmt.Errorf("opcode %02x pushes %d bytes, "+
					"but script only has %d remaining",
					op.Value, int(l), len(script[off:]))
				return retScript, err
			}

			pop.Data = script[off : off+int(l)]
			i += 1 - op.Length + int(l)
		}

		retScript = append(retScript, pop)
	}

	return retScript, nil
}

// RedeemScriptFromTxin extracts Redeem script from the TxIn
func RedeemScriptFromTxin(txin *wire.TxIn) (redeemScript []byte, err error) {

	sigScript, err := parseScript(txin.SignatureScript)
	if err != nil {
		return
	}
	if sigScript == nil {
		err = fmt.Errorf("sigScript parsing error")
		return
	}

	cntr := len(sigScript)
	if cntr < sizeTxInMultisigScriptMin {
		err = fmt.Errorf("unexpected length of TxIn SigScript: %d", cntr)
		return
	}
	cntr--
	// validate initial 0 and the last item
	if sigScript[0].Opcode.Value != OpFalse ||
		sigScript[cntr].Opcode.Value < OpPushData1 ||
		sigScript[cntr].Opcode.Value > OpPushData4 {

		err = fmt.Errorf("unexpected script: %02x, .., %02x",
			sigScript[0].Opcode.Value,
			sigScript[cntr].Opcode.Value)
		return
	}
	// validate intermediate items, which shall be [01, FF]
	for i := 1; i < cntr; i++ {
		var b byte
		if len(sigScript[i].Data) > 0 {
			b = sigScript[i].Data[0]
		}
		if sigScript[i].Opcode.Value != OpData1 ||
			sigScript[i].Opcode.Length != 2 ||
			len(sigScript[i].Data) != 1 ||
			b != 0xff {

			err = fmt.Errorf("unexpected script item[%d]: %02x, %d, %d, %02x",
				i,
				sigScript[i].Opcode.Value,
				sigScript[i].Opcode.Length,
				len(sigScript[i].Data),
				b)
			return
		}
	}
	return sigScript[cntr].Data, nil
}

func xpubFromElectrumEncoded(xpubData []byte) (xkey *hdkeychain.ExtendedKey, path []uint32, err error) {

	checkSum := chainhash.DoubleHashB(xpubData[:sizeEncodedKey])[:4]
	xpub := make([]byte, sizeEncodedKey+sizeCheckSum)

	copy(xpub, xpubData[:sizeEncodedKey])
	copy(xpub[sizeEncodedKey:], checkSum)
	xpubEncoded := base58.Encode(xpub)

	path = make([]uint32, 2)
	ofs := sizeEncodedKey
	path[0] = binary.LittleEndian.Uint32(xpubData[ofs : ofs+4])
	ofs += 4
	path[1] = binary.LittleEndian.Uint32(xpubData[ofs : ofs+4])

	xkey, err = hdkeychain.NewKeyFromString(xpubEncoded)
	if err != nil {
		return
	}
	if xkey == nil {
		err = fmt.Errorf("empty public key")
		return
	}
	/*
	   for _, index := range path {
	       xkey, err = xkey.Child(index)
	       if err != nil {
	       }
	       if xkey == nil {
	           log.Fatalf("empty child key")
	           return
	       }
	   }
	*/
	return
}

// ChildFromXkeyPath returns child of n-th depth from Xpub (n = len(path))
func ChildFromXkeyPath(xkey *hdkeychain.ExtendedKey, path []uint32) (*hdkeychain.ExtendedKey, error) {

	var err error
	for _, index := range path {
		xkey, err = xkey.Child(uint32(index))
		if err != nil {
			return nil, err
		}
		if xkey == nil {
			err = fmt.Errorf("empty child key")
			return nil, err
		}
	}
	return xkey, nil
}

// PubkeysIndexPathFromScript extract xpubs from redeen script,
// searches for the xpub provided and returns it index in the result list,
// builds and returns a list of unsorted pubkeys
func PubkeysIndexPathFromScript(redeem []byte, xpub *ecdsa.PublicKey) (
	m byte, pubkeys []*btcec.PublicKey, index int, xpath []uint32, err error) {

	var pubSample *btcec.PublicKey
	index = -1 // not found

	// parse redeemScript
	msScript, err := parseScript(redeem)
	if err != nil {
		return
	}
	if msScript == nil {
		err = fmt.Errorf("msScript is empty")
		return
	}
	sz := len(msScript)
	if sz < sizeMultisigScriptMin ||
		sz > sizeMultisigScriptMax ||
		msScript[sz-1].Opcode.Value != OpCheckMultiSig {

		err = fmt.Errorf("unexpected script: %d, %02x", sz, msScript[sz-1].Opcode.Value)
		return
	}
	var n byte
	m = msScript[0].Opcode.Value
	n = msScript[sz-2].Opcode.Value
	if m < Op1 || m > Op16 ||
		n < Op1 || n > Op16 ||
		m > n {

		err = fmt.Errorf("unexpected script (%d/%d/%d)", m, n, sz)
		return
	}

	m = m - Op1 + 1
	n = n - Op1 + 1
	if n != byte(sz-sizeMultisigScriptMNCheck) {

		err = fmt.Errorf("unexpected keys number %d, expected %d", n, sz-sizeMultisigScriptMNCheck)
		return
	}

	pubkeys = make([]*btcec.PublicKey, n)
	if xpub != nil {
		pubSample = (*btcec.PublicKey)(xpub)
	}

	ofs := 1
	for i := ofs; i < sz-2; i++ {
		var b byte
		if len(msScript[i].Data) > 0 {
			b = msScript[i].Data[0]
		}

		if msScript[i].Opcode.Value != OpPushData1 ||
			len(msScript[i].Data) != sizeElectrumXpubHD ||
			b != markerXpub {

			err = fmt.Errorf("unexpected xpub data [%d]: %02x, %d, %02x",
				i-ofs,
				msScript[i].Opcode.Value,
				len(msScript[i].Data),
				b)
			return
		}
		var xpub *hdkeychain.ExtendedKey
		var path []uint32
		xpub, path, err = xpubFromElectrumEncoded(msScript[i].Data[1:])
		if err != nil {
			return
		}

		var pub *btcec.PublicKey
		// Extract pubkey from xpub to compare with provided xpub
		if pubSample != nil {
			pub, err = xpub.ECPubKey()
			if err != nil {
				return
			}
			if pubSample.IsEqual(pub) {
				if index != -1 || xpath != nil {
					err = fmt.Errorf("wrong xpub set in script: %d/%d", index, i-ofs)
					return
				}
				index = i - ofs
				xpath = path
			}
		}

		xpub, err = ChildFromXkeyPath(xpub, path)
		if err != nil {
			return
		}

		pub, err = xpub.ECPubKey()
		if err != nil {
			return
		}
		if pub == nil {
			err = fmt.Errorf("empty public key [%d]", i-ofs)
			return
		}
		pubkeys[i-ofs] = pub
	}
	return
}

// MultisigScriptFromPubkeys builds multisig script for TxIn
// returns index of provided pubkey in the sorted list
func MultisigScriptFromPubkeys(m byte, pubkeys []*btcec.PublicKey, index int) (script []byte, indexSorted int, err error) {

	indexSorted = -1
	n := len(pubkeys)
	if n <= 0 {
		err = fmt.Errorf("pubkeys are absent")
		return
	}
	if index < 0 || index >= len(pubkeys) {
		err = fmt.Errorf("invalid index '%d'", index)
		return
	}
	pubkeysCompressed := make([][]byte, len(pubkeys))
	for i := range pubkeys {
		pubkeysCompressed[i] = pubkeys[i].SerializeCompressed()
	}
	searchingPubkey := pubkeysCompressed[index][:]
	sortedPubkeysCompressed := SortedPubkeys(pubkeysCompressed)
	for i := range pubkeys {
		if bytes.Compare(searchingPubkey, sortedPubkeysCompressed[i]) == 0 {
			indexSorted = i
			break
		}
	}
	scriptSize := 1 + (1+btcec.PubKeyBytesLenCompressed)*n + 2
	script = make([]byte, scriptSize)
	script[0] = Op1 + m - 1
	ofs := 1
	for i := range sortedPubkeysCompressed {
		script[ofs] = btcec.PubKeyBytesLenCompressed
		ofs++
		copy(script[ofs:], sortedPubkeysCompressed[i])
		ofs += btcec.PubKeyBytesLenCompressed
	}
	script[scriptSize-2] = Op1 + byte(n) - 1
	script[scriptSize-1] = OpCheckMultiSig
	return
}
