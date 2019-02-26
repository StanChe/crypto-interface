package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/stanche/crypto-interface/connector"
	"github.com/stanche/crypto-interface/connector/btc/script"

	"github.com/wedancedalot/decimal"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	btcchaincfg "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

const (
	scriptVerifyFailureErr = "mandatory-script-verify-flag-failed"

	defaultTimeoutSec = 30
)

type (
	AddressDecoder func(string) (btcutil.Address, error)

	IBtcChainConnector interface {
		connector.IConnector
		GetBlockByNumber(number uint64) (*wire.MsgBlock, error)
		GetTransactionByHash(hash chainhash.Hash) (tx *btcjson.TxRawResult, isPending bool, err error)
		ParseOutputs(txOut []*wire.TxOut) ([]*connector.OutputParsed, error)
		CreateRawTransaction(inputs []btcjson.TransactionInput,
			amounts map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error)
		DecoderSet(decoder AddressDecoder)
	}

	BtcChainConnector struct {
		connector.Connector
		NodeURL    string
		Client     *rpcclient.Client
		chain      *chaincfg.Params
		Decoder    AddressDecoder
		CoreClient *Client
	}
)

func clientUrl(cfg connector.NodeParams) (string, error) {
	if cfg.Host != "" && cfg.Port != 0 && cfg.User != "" && cfg.Password != "" {
		return fmt.Sprintf("http://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port), nil
	}
	return "", fmt.Errorf("invalid config")
}

func NewBtcChainConnector(walletID uint64, cfg *connector.WalletParams) (IBtcChainConnector, error) {

	var err error
	if cfg == nil || walletID <= 0 {
		err = fmt.Errorf("Wallet configuration parameters absent")
		// log.Errorf(err.Error())
		return nil, err
	}
	if !cfg.Active {
		return &BtcChainConnector{}, nil
	}
	connector := &BtcChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   cfg.Currency,
			WalletType: cfg.Type,
		},
		chain: &btcchaincfg.TestNet3Params,
	}
	connector.DecoderSet(connector.DecodeAddress)

	connector.NodeURL, err = clientUrl(cfg.Node)
	if err != nil {
		return nil, err
	}
	connector.Client, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", cfg.Node.Host, cfg.Node.Port),
		User:         cfg.Node.User,
		Pass:         cfg.Node.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		// log.Errorf("failed to connect to Bitcoin node: %v", err.Error())
	}
	coreURL, err := clientUrl(cfg.Core)
	if err == nil {
		timeout := defaultTimeoutSec
		if cfg.NodeTimeoutSec > 0 {
			timeout = cfg.NodeTimeoutSec
		}
		connector.CoreClient = NewClient(coreURL, timeout)
	} else {
		// log.Errorf("bitcore clientUrl error:%s", err.Error())
	}
	return connector, nil
}

const btcPrecision = 8

type coreBalance struct {
	Balance int64 `json:"balance"`
	// received ...
}

func (c *BtcChainConnector) balance(addr string) (int64, error) {

	if c.CoreClient == nil || c.CoreClient.URL == "" {
		return 0, fmt.Errorf("coreClient not initialized")
	}
	data := fmt.Sprintf(`{"jsonrpc": "1.0", "id":"core", "method": "getaddressbalance", "params": ["%s"] }`, addr)
	resp, err := c.CoreClient.send(data)
	if err != nil {
		return 0, err
	}
	var res coreBalance
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return 0, err
	}
	return res.Balance, nil
}

func (c *BtcChainConnector) BalanceGet(currency connector.Currency, addresses ...string) (b connector.AddressBalance, err error) {

	if len(addresses) == 0 {
		// log.Errorf("btcChainConnector does not support BalanceGet with empty addresses list")
		return b, fmt.Errorf("unsupported params: BalanceGet.addresses are empty")
	}
	total := big.NewInt(0)
	var valid bool
	var balance int64
	for _, addr := range addresses {
		valid, err = c.ValidateAddress(addr)
		if err != nil {
			// log.Errorf("%s.ChainConnector.BalanceGet.ValidateAddress: %s", c.CurrencyCode(), err.Error())
			return b, err
		}
		if !valid {
			continue
		}
		balance, err = c.balance(addr)
		if err != nil {
			// log.Errorf("balance(%s): %s", addr, err.Error())
			return b, err
		}
		//total = total.Add(total, balance)
		total = total.Add(total, big.NewInt(balance))
	}
	return connector.AddressBalance{
		Confirmed: decimal.NewFromBigInt(total, 0).Div(decimal.New(1, int32(currency.GetPrecision()))),
	}, nil

}

func (c *BtcChainConnector) ValidateAddress(address string) (bool, error) {
	addr, err := btcutil.DecodeAddress(address, c.chain)
	if err != nil {
		// log.Errorf("ValidateAddress[%s] (%s): %s", c.Currency, address, err.Error())
		return false, nil
	}
	return addr.IsForNet(c.chain), nil
}

func (c *BtcChainConnector) DecodeAddress(addr string) (btcutil.Address, error) {
	return btcutil.DecodeAddress(addr, c.chain)
}

func (c *BtcChainConnector) DecoderSet(decoder AddressDecoder) {
	c.Decoder = decoder
}

func (c *BtcChainConnector) ParseOutputs(txOuts []*wire.TxOut) ([]*connector.OutputParsed, error) {
	var outputs []*connector.OutputParsed
	for i, txOut := range txOuts {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, c.chain)
		if err != nil {
			// log.Errorf("ExtractTxOutAddresses %s", err.Error())
			//todo find out what to do in case if we cnnot parse output address
			continue
		}

		for _, address := range addresses {
			outputs = append(outputs, &connector.OutputParsed{
				Address: address.EncodeAddress(),
				Value:   big.NewInt(txOut.Value),
				TxPos:   uint(i),
			})
		}
	}
	return outputs, nil
}

func (c *BtcChainConnector) CreateRawTransaction(inputs []btcjson.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {

	return c.Client.CreateRawTransaction(inputs, amounts, nil)
}

func (c *BtcChainConnector) GetBlockByNumber(number uint64) (*wire.MsgBlock, error) {
	blockHash, err := c.Client.GetBlockHash(int64(number))
	if blockHash == nil || err != nil {
		return nil, connector.ErrNotFound
	}

	block, err := c.Client.GetBlock(blockHash)
	if block == nil || err != nil {
		return nil, connector.ErrNotFound
	}
	return block, nil
}
func (c *BtcChainConnector) GetTransactionByHash(hash chainhash.Hash) (*btcjson.TxRawResult, bool, error) {
	tx, err := c.Client.GetRawTransactionVerbose(&hash)

	pending := false
	if err == nil && tx != nil {
		pending = tx.BlockHash == ""
	}
	return tx, pending, err
}

// TxStatus returns transaction status by TxId(hash)
func (c *BtcChainConnector) TxStatus(txID string, blockNo uint64) (*connector.TxStatusStruct, error) {

	txHash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, fmt.Errorf("btc wallet Client NewHashFromStr failed: %s", err.Error())
	}

	tx, err := c.Client.GetRawTransactionVerbose(txHash)
	if err != nil {
		return nil, fmt.Errorf("btc wallet Client GetRawTransactionVerbose failed: %s", err.Error())
	}

	if tx == nil {
		return nil, nil
	}

	if tx.BlockHash == "" {
		//Unconfirmed
		status := connector.TxStatusStruct{
			Conf:   0,
			Height: 0,
		}
		return &status, nil
	}
	status := connector.NewTxStatusWithNonNeg(0, int64(tx.Confirmations))
	return &status, nil
}

func (c *BtcChainConnector) TxBuild(walletData *connector.WalletSignStruct,
	utxosIn interface{}, output []connector.OutStruct) (string, error) {

	n := len(walletData.XPubs)
	if n < 2 || n > 15 {
		return "", fmt.Errorf("invalid signers quantity")
	}
	m := int(walletData.Signers)
	if m > n || m < 1 || m > 15 {
		return "", fmt.Errorf("invalid signers required number")
	}

	utxos, ok := utxosIn.([]connector.UtxStruct)
	if !ok {
		return "", fmt.Errorf("unexpected type of utxo input for BTC wallet: expected []connector.UtxStruct got: %+v", utxosIn)
	}

	inputs := make([]btcjson.TransactionInput, len(utxos))
	for i := range utxos {
		inputs[i] = btcjson.TransactionInput{
			Txid: utxos[i].TxHash,
			Vout: uint32(utxos[i].TxPos),
		}
	}

	values := make(map[string]decimal.Decimal)
	var val, value decimal.Decimal
	for i := range output {
		value = output[i].Amount.Abs() // just to copy value
		val, ok = values[output[i].Address]
		if ok {
			value = value.Add(val)
		}
		values[output[i].Address] = value
	}

	amounts := make(map[btcutil.Address]btcutil.Amount)
	var address btcutil.Address
	var err error
	for addr, amt := range values {
		address, err = c.Decoder(addr)
		if err != nil {
			return "", err
		}
		amounts[address] = btcutil.Amount(amt.Mul(decimal.New(1, int32(btcPrecision))).IntPart())
	}

	msg, err := c.CreateRawTransaction(inputs, amounts)
	if err != nil {
		return "", err
	}

	for inputNo := range msg.TxIn {
		err = ScriptBuild(msg.TxIn[inputNo], utxos[inputNo].Index, int(walletData.Signers), walletData.XPubs, nil)
		if err != nil {
			return "", err
		}
	}

	var b bytes.Buffer
	b.Grow(msg.SerializeSize())
	err = msg.Serialize(&b)
	if err != nil {
		return "", err
	}
	//fmt.Printf("hexTx: %s", hex.EncodeToString(b.Bytes()))
	return hex.EncodeToString(b.Bytes()), nil
}

func ScriptBuild(txIn *wire.TxIn, index uint32,
	signaturesRequired int, xpubs []string, signatures []string) error {

	const xpubSize = 1 + 4 + 1 + 4 + 4 + 32 + 1 + 32 + 4 + 4
	// ff 0488b21e 00 00000000 00000000
	// d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c6
	// 03 d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a59
	// 00000000 00000000
	const xpubValueSize = 4 + 1 + 4 + 4 + 32 + 1 + 32
	const xpubIndexOffset = 1 + 4 + 1 + 4 + 4 + 32 + 1 + 32 + 4
	const pubSize = 1 + 32
	const pubValueSize = 32

	var err error

	walletM := signaturesRequired
	walletN := len(xpubs)

	flagSignatures := len(signatures) > 0 && walletM == len(signatures)

	var bufSize int
	if flagSignatures {
		bufSize = pubSize
	} else {
		bufSize = xpubSize
	}

	indexBuf := make([]byte, 4)
	offsets := make([]int, walletN)

	builder := txscript.NewScriptBuilder()
	builder.AddOp(byte(txscript.OP_1 + walletM - 1))
	for i := 0; i < walletN; i++ {
		xpub := make([]byte, bufSize)
		if flagSignatures {
			xpub, err = hex.DecodeString(xpubs[i])
			if err != nil {
				return err
			}
		} else {
			xpubDecoded := base58.Decode(xpubs[i])
			xpub[0] = 0xff
			copy(xpub[1:1+xpubValueSize], xpubDecoded[0:xpubValueSize])
		}

		builder.AddData(xpub)
		if !flagSignatures {
			scriptBuf, err := builder.Script()
			if err != nil {
				return err
			}
			offsets[i] = len(scriptBuf) - 4
		}
	}
	builder.AddOp(byte(txscript.OP_1 + walletN - 1))
	builder.AddOp(txscript.OP_CHECKMULTISIG)
	pkScript, err := builder.Script()
	if err != nil {
		return err
	}

	scriptBuilder := txscript.NewScriptBuilder()
	scriptBuilder.AddOp(txscript.OP_0)
	// fill signatures
	for j := 0; j < walletM; j++ {
		var signature []byte
		if flagSignatures {
			signature, err = hex.DecodeString(signatures[j])
			if err != nil {
				return err
			}
		} else {
			signature = []byte{0xff}
		}
		scriptBuilder.AddData(signature)
	}
	if !flagSignatures {
		// update public key indexes
		for j := 0; j < walletN; j++ {
			// put index into indexBuf
			binary.LittleEndian.PutUint32(indexBuf, index)
			// substitute index value in the pkScript
			copy(pkScript[offsets[j]:], indexBuf[:])
		}
	}
	// append [modified] pkScript
	scriptBuilder.AddData(pkScript)
	script, err := scriptBuilder.Script()
	if err != nil {
		return err
	}
	txIn.SignatureScript = script[:]
	return nil
}

func (c *BtcChainConnector) TxBroadcast(txHex string) (string, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return "", err
	}
	txReader := bytes.NewBuffer(txBytes)

	var msg wire.MsgTx
	err = msg.Deserialize(txReader)
	if err != nil {
		return "", err
	}

	hash, err := c.Client.SendRawTransaction(&msg, false)
	if err != nil {
		// log.Errorf("SendRawTransaction failed for %s: %s", c.CurrencyCode(), err.Error())
		if strings.Contains(err.Error(), scriptVerifyFailureErr) {
			return "", connector.TxPermanentFailure
		}
		return "", err
	}
	return hash.String(), nil
}

// TxRebuild - combine parsed hex Tx with the signatures
func (c *BtcChainConnector) TxRebuild(txHex string, signatures connector.TxSignatures) (string, error) {
	return TxRebuildBtc(txHex, signatures)
}

func TxRebuildBtc(txHex string, signatures connector.TxSignatures) (string, error) {

	txData, err := hex.DecodeString(txHex)
	if err != nil {
		return "", err
	}

	tx, err := btcutil.NewTxFromBytes(txData)
	if err != nil {
		// log.Errorf("SignProcessor[BTC].NewTxFromBytes: %s", err.Error())
		return "", err
	}
	if tx == nil {
		err = fmt.Errorf("decoded tx is nil")
		// log.Errorf("SignProcessor[BTC].NewTxFromBytes: %s", err.Error())
		return "", err
	}
	msgTx := tx.MsgTx()
	countTxIn := len(msgTx.TxIn)
	if countTxIn != len(signatures) {
		err = fmt.Errorf("inconsistent tx inputs and signatures quantity: %d ~ %d", countTxIn, len(signatures))
		return "", err
	}
	for indexTxIn := 0; indexTxIn < countTxIn; indexTxIn++ {

		redeemScript, err := script.RedeemScriptFromTxin(tx.MsgTx().TxIn[indexTxIn])
		if err != nil {
			return "", err
		}

		m, pubkeys, _, xpath, err := script.PubkeysIndexPathFromScript(redeemScript, nil)
		if err != nil {
			return "", err
		}

		if len(signatures[indexTxIn]) != int(m) {
			return "", fmt.Errorf("inconsistent signatures (%d, expected %d) for input %d",
				len(signatures[indexTxIn]), m, indexTxIn)
		}

		if len(xpath) != 2 && xpath[0] != uint32(0) {
			return "", fmt.Errorf("invalid xpath calculated for input %d", indexTxIn)
		}

		pubs := make([]string, len(pubkeys))
		for i := range pubkeys {
			pubs[i] = hex.EncodeToString(pubkeys[i].SerializeCompressed())
		}
		sort.Strings(pubs)
		err = ScriptBuild(msgTx.TxIn[indexTxIn], xpath[1], int(m), pubs, signatures[indexTxIn])
		if err != nil {
			return "", err
		}
	}

	var b bytes.Buffer
	b.Grow(msgTx.SerializeSize())
	err = msgTx.Serialize(&b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}
