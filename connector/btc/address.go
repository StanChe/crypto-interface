package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
)

type (
	// CoinAddress is a simplified type which implements btcutil.Address interface
	CoinAddress struct {
		Addr string
	}
)

// String returns the string encoding of the transaction output
// destination.
//
// Please note that String differs subtly from EncodeAddress: String
// will return the value as a string without any conversion, while
// EncodeAddress may convert destination types (for example,
// converting pubkeys to P2PKH addresses) before encoding as a
// payment address string.
func (a *CoinAddress) String() string {
	return a.Addr
}

// EncodeAddress returns the string encoding of the payment address
// associated with the Address value.  See the comment on String
// for how this method differs from String.
func (a *CoinAddress) EncodeAddress() string {
	return a.Addr
}

// ScriptAddress returns the raw bytes of the address to be used
// when inserting the address into a txout's script.
func (a *CoinAddress) ScriptAddress() []byte {
	scriptAddr, _, err := base58.CheckDecode(a.Addr)
	if err != nil {
		return nil
	}
	return scriptAddr
}

// IsForNet returns whether or not the address is associated with the
// passed bitcoin network.
func (a *CoinAddress) IsForNet(*chaincfg.Params) bool {
	return false
}
