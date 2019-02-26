package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stanche/crypto-interface/connector"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/wedancedalot/decimal"
)

type Currency struct {
	Code         string
	Precision    uint8
	TokenAddress string
	TokenCode    int64
}

// GetTokenAddress gets the address of the contract of the currency (if applicable). If not applicable - returns an emppty string.
func (c Currency) GetTokenAddress() string {
	return c.TokenAddress
}

// GetCode is the code of the currency (i.e. BTC/ETH/USDT). It's usually capitalized.
func (c Currency) GetCode() string {
	return c.Code
}

// GetPrecision gets the maximum number of decimal points of the currency.
func (c Currency) GetPrecision() uint8 {
	return c.Precision
}

// GetTokenCode gets an integer code of the currency (if applicable). If not applicable - returns 0.
func (c Currency) GetTokenCode() int64 {
	return c.TokenCode
}

func Test_nodeConnector_TxBuild(t *testing.T) {
	t.Skip("BTC Node IConnector integration test is skipped") //comment this line to run tests

	type fields struct {
		client *rpcclient.Client
		chain  *chaincfg.Params
	}
	type args struct {
		walletData *connector.WalletSignStruct
		utxosIn    interface{}
		output     []connector.OutStruct
	}

	amount, _ := decimal.NewFromString("0.665983")
	amountW1, _ := decimal.NewFromString("0.09")
	amountW2, _ := decimal.NewFromString("0.09")
	amountW3, _ := decimal.NewFromString("9.99")
	amountW4, _ := decimal.NewFromString("0.99")

	c := connector.NodeParams{
		Host:     "104.199.25.196",
		Port:     8341,
		User:     "attic",
		Password: "8Wsujmq4JND0565itTqt",
	}
	clientDEV, _ := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:         c.User,
		Pass:         c.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{

		{
			name: "dummy multysig tx",
			fields: fields{
				client: clientDEV,
				chain:  nil,
			},

			args: args{
				walletData: &connector.WalletSignStruct{
					Signers: uint8(2),
					XPubs: []string{
						"xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufqeGBwALkBUK",
						"xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ",
						"xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv",
					},
				},
				utxosIn: []connector.UtxStruct{
					connector.UtxStruct{
						TxHash: "8fc527b744bf0784bfeb8b610386e3e437ca4796a5ce9bf7c728bebe622717db",
						TxPos:  0,
						Index:  uint32(2),
					},
				},
				output: []connector.OutStruct{
					connector.OutStruct{
						Address:  "2MtBe9ZJwGV8eJDdJkytbuq8y5gwB9HxxC3",
						Amount:   amount,
						Currency: &Currency{Code: "BTC"},
						Memo:     "",
					},
				},
			},
			want:    "0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000",
			wantErr: false,
		},

		{
			name: "dummy multysig tx",
			fields: fields{
				client: clientDEV,
				chain:  nil,
			},

			args: args{
				walletData: &connector.WalletSignStruct{
					Signers: uint8(2),
					XPubs: []string{
						"xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufqeGBwALkBUK",
						"xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ",
						"xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv",
					},
				},
				utxosIn: []connector.UtxStruct{
					connector.UtxStruct{
						TxHash: "8fc527b744bf0784bfeb8b610386e3e437ca4796a5ce9bf7c728bebe622717db",
						TxPos:  0,
						Index:  uint32(2),
					},
				},
				output: []connector.OutStruct{
					connector.OutStruct{
						Address: "2MtBe9ZJwGV8eJDdJkytbuq8y5gwB9HxxC3",
						Amount:  amount,
					},
					connector.OutStruct{
						Address: "n12fkNBS9XuQXRscN1k62xaK1r6pT215cW",
						Amount:  amountW1,
					},
					connector.OutStruct{
						Address: "n12fkNBS9XuQXRscN1k62xaK1r6pT215cW",
						Amount:  amountW2,
					},
					connector.OutStruct{
						Address: "n12fkNBS9XuQXRscN1k62xaK1r6pT215cW",
						Amount:  amountW3,
					},
					connector.OutStruct{
						Address: "n12fkNBS9XuQXRscN1k62xaK1r6pT215cW",
						Amount:  amountW4,
					},
				},
			},
			want:    "0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &BtcChainConnector{
				Client: tt.fields.client,
				chain:  tt.fields.chain,
			}
			nc.DecoderSet(nc.DecodeAddress)
			got, err := nc.TxBuild(tt.args.walletData, tt.args.utxosIn, tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeConnector.TxBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nodeConnector.TxBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeConnector_TxRebuild(t *testing.T) {
	type fields struct {
		client *rpcclient.Client
		chain  *chaincfg.Params
	}
	type args struct {
		txHex      string
		signatures connector.TxSignatures
	}

	c := connector.NodeParams{
		Host:     "104.199.25.196",
		Port:     8341,
		User:     "attic",
		Password: "8Wsujmq4JND0565itTqt",
	}
	clientDEV, _ := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:         c.User,
		Pass:         c.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "multysig 2 of 3 rebuild",
			fields: fields{
				client: clientDEV,
				chain:  nil,
			},

			args: args{
				txHex: "0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000",
				signatures: []connector.TxInSignatures{
					connector.TxInSignatures{
						"30440220596c276e66186b98e1b190a626a94b30760718c99f2db32d2e165e7075c3f67302207fd6cd72995239952769b7ff2c61f4e952a05a3a9970f564f09f3efe51feded201",
						"3045022100cc08a8be0f1021f9029b0fd428a0d1575c39e215ee396672eb70dd350f5b17d30220075c8eaf70f6d0dfcdf379c4a7f98ae7eb64ef765ac8d7be1b30f6e7e4c4181301",
					},
				},
			},
			want:    "0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fdfd00004730440220596c276e66186b98e1b190a626a94b30760718c99f2db32d2e165e7075c3f67302207fd6cd72995239952769b7ff2c61f4e952a05a3a9970f564f09f3efe51feded201483045022100cc08a8be0f1021f9029b0fd428a0d1575c39e215ee396672eb70dd350f5b17d30220075c8eaf70f6d0dfcdf379c4a7f98ae7eb64ef765ac8d7be1b30f6e7e4c41813014c69522102e9686c62273b60cdf58ee4c8bda595780bcdf5b441161b7301dad939dbe83ec42103297c46de43997b7f6702d9c353ffb2566d818eed8c295b16ed46076f328474882103839cd7a8f5fe9ef2528325425761d9c78bb162f7a70b8f4d88235cf763ee13ca53aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &BtcChainConnector{
				Client: tt.fields.client,
				chain:  tt.fields.chain,
			}
			got, err := nc.TxRebuild(tt.args.txHex, tt.args.signatures)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeConnector.TxRebuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nodeConnector.TxRebuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBtcConnector_IntegrationTest(t *testing.T) {

	t.Skip("BTC Node IConnector integration test is skipped") //comment this line to run tests

	//Given
	nodeConfig := &connector.NodeParams{
		Host:     "104.199.25.196",
		Port:     8361,
		User:     "attic",
		Password: "8Wsujmq4JND0565itTqt",
	}

	bc, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", nodeConfig.Host, nodeConfig.Port),
		User:         nodeConfig.User,
		Pass:         nodeConfig.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}

	params := &chaincfg.MainNetParams
	walletConnector := &BtcChainConnector{
		Client: bc,
		chain:  params,
	}
	/*
	   	//block, err := walletConnector.GetBlockByNumber(uint64(539187))
	   	block, err := walletConnector.GetBlockByNumber(uint64(214696))
	   	if err != nil {
	   		// log.Fatalf("walletConnector.GetBlockByNumber: %v", err.Error())
	   	}
	   	// log.Debugf("Block %s", block.BlockHash())

	   	parsedOutputs, err := walletConnector.ParseOutputs(block.Transactions[0].TxOut)
	   	if err != nil {
	   		// log.Fatalf("walletConnector.ParseOutputs: %v", err.Error())
	   	}
	       // log.Debugf("parsedOutputs %v", parsedOutputs)
	*/
	balance, err := walletConnector.BalanceGet(Currency{}, "2MyhnviNUUMrZke1uhGwhAxmUYvvGb14MbA")
	if err != nil {
		//// log.Fatalf("walletConnector.ParseOutputs: %v", err.Error())
	}
	fmt.Printf("hot wallet balance: %v", balance)
}

func Test_Decode_Encode_OMNI_tx(t *testing.T) {
	omniTxHex := "0100000002bf66dbc91cc28dfc9f60343a1ac9d473ca0d4949190883265b66e9c454a4911d0200000000ffffffff8c5ef5be9969c51b0afe621d392a36246bd8fd1ac4d150c762ddb96614248a3c0000000000ffffffff034ab4eb0b000000001976a914970ba8eb3925e699b4a5b6c56ec70820c154eb7488ac1c0200000000000017a9148f9e7981344e568210e7926aaf729b25e369a18f870000000000000000166a146f6d6e6900000000000000020000000005f5e10000000000"

	txData, err := hex.DecodeString(omniTxHex)
	if err != nil {
		// log.Fatalf("DecodeString: %s", err.Error())
	}

	tx, err := btcutil.NewTxFromBytes(txData)
	if err != nil {
		// log.Fatalf("NewTxFromBytes: %s", err.Error())
	}

	msgTx := tx.MsgTx()

	var b bytes.Buffer
	b.Grow(msgTx.SerializeSize())
	err = msgTx.Serialize(&b)
	if err != nil {
		// log.Fatalf("Serialize: %s", err.Error())
	}
	txHex := hex.EncodeToString(b.Bytes())
	// log.Debugf("omniTxHex: %s", omniTxHex)
	// log.Debugf("hexTx: %s", txHex)

	assert.Equal(t, omniTxHex, txHex, "Unexpected value")
}

func TestBtcConnector_ValidateAddress(t *testing.T) {
	//Given
	nodeConfig := &connector.NodeParams{
		Host:     "104.199.25.196",
		Port:     8331,
		User:     "attic",
		Password: "8Wsujmq4JND0565itTqt",
	}

	bc, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", nodeConfig.Host, nodeConfig.Port),
		User:         nodeConfig.User,
		Pass:         nodeConfig.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}

	addressMainNet1 := "39hxwP9ZGm1DcKW3hDwWjWwdQVcvYhuLnd"
	addressMainNet2 := "34ZqX3zRf1br4qFGjUgv9RjSKeAhLdDtrh"

	//params :=
	//params := &chaincfg.TestNet3Params
	walletConnectorMainNet := &BtcChainConnector{
		Client: bc,
		chain:  &chaincfg.MainNetParams,
	}
	walletConnectorTestNet := &BtcChainConnector{
		Client: bc,
		chain:  &chaincfg.TestNet3Params,
	}

	validationResult1, err := walletConnectorMainNet.ValidateAddress(addressMainNet1)
	// log.Debugf("validationResult(MainNet) %v", validationResult1)

	validationResult2, err := walletConnectorTestNet.ValidateAddress(addressMainNet2)
	// log.Debugf("validationResult(TestNet) %v", validationResult2)

	assert.Equal(t, true, validationResult1, "Unexpected value")
	assert.Equal(t, false, validationResult2, "Unexpected value")
}

func TestBtcChainConnector_balance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Skip("BTC Node IConnector integration test is skipped") //comment this line to run tests
		// port 25012 from cryptagio-coins-testnet-0 shall be mapped into local

		nodeConfig := connector.NodeParams{
			Host:     "127.0.0.1",
			Port:     23002,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		}
		coreConfig := connector.NodeParams{
			Host:     "127.0.0.1",
			Port:     25012,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		}
		config := &connector.WalletParams{
			Active:         true,
			Node:           nodeConfig,
			NodeTimeoutSec: 30,
			Core:           coreConfig,
		}

		currency := Currency{}
		walletID := uint64(1)
		conn, _ := NewBtcChainConnector(walletID, config)

		addressDev := "mqpaRTpgKSnbeaqWmS9cwEobjtFVHsmjuX"
		expectedBalance := int64(299994314)
		expTotal := decimal.NewFromBigInt(big.NewInt(expectedBalance), 0)

		devBalance, err := conn.BalanceGet(currency, addressDev)
		assert.Nilf(t, err, "unexpected error")
		cmp := expTotal.Cmp(devBalance.Confirmed)
		assert.Equal(t, 0, cmp, "unexpected value")
	})
	t.Run("it should return error if core client not set", func(t *testing.T) {
		nodeConfig := connector.NodeParams{
			Host:     "127.0.0.1",
			Port:     23002,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		}
		currency := Currency{}
		walletID := uint64(1)
		config := &connector.WalletParams{
			Active:         true,
			Node:           nodeConfig,
			NodeTimeoutSec: 30,
		}
		addressDev := "mqpaRTpgKSnbeaqWmS9cwEobjtFVHsmjuX"
		conn, _ := NewBtcChainConnector(walletID, config)
		_, err := conn.BalanceGet(currency, addressDev)
		assert.NotNil(t, err, "expect error")
		assert.Containsf(t, err.Error(), "coreClient not initialized", "should contain error message about core client")
	})
}
