package connector

import (
	"fmt"
)

type (
	// TxInSignatures - signatures array for a tx input
	TxInSignatures []string
	// TxSignatures - signatures set for the tx
	TxSignatures []TxInSignatures
	// Currency provides info about the currency.
	Currency interface {
		// GetCode is the code of the currency (i.e. BTC/ETH/USDT). It's usually capitalized.
		GetCode() string
		// GetPrecision gets the maximum number of decimal points of the currency.
		GetPrecision() uint8
		// GetTokenAddress gets the address of the contract of the currency (if applicable). If not applicable - returns an emppty string.
		GetTokenAddress() string
		// GetTokenCode gets an integer code of the currency (if applicable). If not applicable - returns 0.
		GetTokenCode() int64
	}

	// Address provides info about the wallet address
	Address interface {
		// GetAddress is the address itself
		GetAddress() string
		// GetTag returns the tag or memo property for currencies, that uses it
		GetTag() string
	}

	// BalanceProvider is an interface for getting sum of the balances on specified addresses.
	BalanceProvider interface {
		BalanceGet(currency Currency, address ...string) (AddressBalance, error)
	}
	// AddressValidator is an interface for verifying whether the address is valid for the blockchain.
	// If it returns true - we are safe to send coins to that address.
	AddressValidator interface {
		ValidateAddress(address string) (bool, error)
	}
	// TxBuilder is an interface for building transactions.
	TxBuilder interface {
		TxBuild(walletData *WalletSignStruct, utxos interface{}, output []OutStruct) (string, error)
		// TxRebuild combines raw txHex (built with TxBuild) with signatures from the signer
		// and produces a transaction with signatures that is ready for broadcasting.
		TxRebuild(txHex string, signatures TxSignatures) (string, error)
	}

	// TxGetter is an interface for getting tx status from chain.
	TxGetter interface {
		TxStatus(txID string, blockNo uint64) (*TxStatusStruct, error)
	}

	// TxSender is an interface for transaction broadcasting.
	TxSender interface {
		TxBroadcast(txHex string) (txHash string, err error)
	}

	// IConnector defines wallet (node) interface
	IConnector interface {
		TxGetter
		TxBuilder
		TxSender
		BalanceProvider
		AddressValidator
		CurrencyCode() string
		WalletID() uint64
		GetWalletType() string
	}
)

// TxPermanentFailure indicates a permanent failure for a tx - their is no need to try to broadcast the tx that returns such error.
var TxPermanentFailure = fmt.Errorf("transaction failed permanently")

var (
	ErrNotFound  = fmt.Errorf("not found")
	ErrClientNil = fmt.Errorf("client is nil")
)
