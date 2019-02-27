package bch

import (
	"cryptagio-walle/dao/models"
	"fmt"
	"testing"

	bchchaincfg "github.com/bchsuite/bchd/chaincfg"
	"github.com/stanche/crypto-interface/connector"
	"github.com/stanche/crypto-interface/connector/btc_example"

	"github.com/stretchr/testify/assert"
	"github.com/wedancedalot/decimal"
)

func Test_nodeConnector_TxBuild(t *testing.T) {

	t.Skip("BCHABC Node IConnector integration test is skipped") //comment this line to run tests

	type args struct {
		walletData *connector.WalletSignStruct
		utxosIn    interface{}
		output     []connector.OutStruct
	}

	amount, _ := decimal.NewFromString("0.01")
	change, _ := decimal.NewFromString("9.999")

	amount1, _ := decimal.NewFromString("1.0")
	change1, _ := decimal.NewFromString("9.099")

	walletID := uint64(1)
	walletConfig := &connector.WalletParams{
		Currency: "BCHABC",
		Type:     "hot",
		Node: connector.NodeParams{
			Host:     "104.199.25.196",
			Port:     8341,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		},
		ChainConfig: "regtest",
	}
	iBtcConnector, err := btc_example.NewBtcChainConnector(walletID, walletConfig)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}
	walletConnector := &bchChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   walletConfig.Currency,
			WalletType: walletConfig.Type,
		},
		ibtc:    iBtcConnector,
		chain:   &bchchaincfg.TestNet3Params,
		regtest: walletConfig.ChainConfig == "regtest",
	}
	walletConnector.ibtc.DecoderSet(walletConnector.DecodeAddress)

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// f6efb8592325c58e9ddc9019bfb09e19bb244beac1d5a774d310f2b114051fd7:0 10.1
		{
			name: "integration multysig tx",
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
						TxHash: "f6efb8592325c58e9ddc9019bfb09e19bb244beac1d5a774d310f2b114051fd7",
						TxPos:  0,
						Index:  uint32(1000),
					},
				},
				output: []connector.OutStruct{
					connector.OutStruct{
						Address:  "bchreg:qzu7d73hakh39hc2qqmz2ln73x56hdp04cyfy5q4ye",
						Amount:   amount1,
						Currency: &models.Currency{Code: string("BCHABC")},
						Memo:     "",
					},
					connector.OutStruct{
						Address:  "bchreg:pzhhpwatsrakfklepvsj7jt3e3yq059csqw8u05deg",
						Amount:   change1,
						Currency: &models.Currency{Code: string("BCHABC")},
						Memo:     "",
					},
				},
			},
			want:    "0200000001d71f0514b1f210d374a7d5c1ea4b24bb199eb0bf1990dc9d8ec5252359b8eff600000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000e80300004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000e80300004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a00000000e803000053aeffffffff02e0f83b360000000017a914af70bbab80fb64dbf90b212f4971cc4807d0b8808700e1f505000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000",
			wantErr: false,
		},
		{
			name: "dummy multysig tx",
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
						TxHash: "b5fbac128e00fd45968468b90fde985b5132d458a66292a39bc4ad4639d0a56b",
						TxPos:  0,
						Index:  uint32(0),
					},
				},
				output: []connector.OutStruct{
					connector.OutStruct{
						Address:  "bchreg:qpwl0qvgkv0qzdh55mr254c0t3q64x2t355a20k24d",
						Amount:   amount,
						Currency: &models.Currency{Code: string("BCHABC")},
						Memo:     "",
					},
					connector.OutStruct{
						Address:  "bchreg:qzu7d73hakh39hc2qqmz2ln73x56hdp04cyfy5q4ye",
						Amount:   change,
						Currency: &models.Currency{Code: string("BCHABC")},
						Memo:     "",
					},
				},
			},
			want:    "02000000016ba5d03946adc49ba39262a658d432515b98de0fb968849645fd008e12acfbb500000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000000000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000000000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000000000053aeffffffff0240420f00000000001976a9145df78188b31e0136f4a6c6aa570f5c41aa994b8d88ac6043993b000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := walletConnector.TxBuild(tt.args.walletData, tt.args.utxosIn, tt.args.output)
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

func TestBchConnector_IntegrationTest(t *testing.T) {

	t.Skip("BCHABC Node IConnector integration test is skipped") //comment this line to run tests

	//Given
	walletID := uint64(1)
	walletConfig := &connector.WalletParams{
		Currency: "BCHABC",
		Type:     "hot",
		Node: connector.NodeParams{
			Host:     "104.199.25.196",
			Port:     8341,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		},
		ChainConfig: "regtest",
	}
	iBtcConnector, err := btc_example.NewBtcChainConnector(walletID, walletConfig)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}
	walletConnector := &bchChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   walletConfig.Currency,
			WalletType: walletConfig.Type,
		},
		ibtc:    iBtcConnector,
		chain:   &bchchaincfg.TestNet3Params,
		regtest: walletConfig.ChainConfig == "regtest",
	}

	txStatus, err := walletConnector.TxStatus("ee215acf6b24a26aa160029a74a3b6ef8ae8984c25d784211d829df16a9dd3c3", 0)
	if err != nil {
		// log.Fatalf("walletConnector.GetBlockByNumber: %v", err.Error())
	}
	fmt.Printf("txStaus %v", txStatus)

	block, err := walletConnector.GetBlockByNumber(uint64(237397))
	if err != nil {
		// log.Fatalf("walletConnector.GetBlockByNumber: %v", err.Error())
	}
	fmt.Printf("Block %s", block.BlockHash())

	for txNo := range block.Transactions {
		parsedOutputs, err := walletConnector.ParseOutputs(block.Transactions[txNo].TxOut)
		if err != nil {
			// log.Fatalf("walletConnector.ParseOutputs: %v", err.Error())
		}
		for i := range parsedOutputs {
			fmt.Printf("parsedOutput[%d]: %v", i, parsedOutputs[i])
		}
	}
	fmt.Printf("done")
}

func TestBchConnector_ValidateAddress(t *testing.T) {
	//Given
	walletID := uint64(1)
	walletConfig := &connector.WalletParams{
		Currency: "BCHABC",
		Type:     "hot",
		Node: connector.NodeParams{
			Host:     "104.199.25.196",
			Port:     8341,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		},
		ChainConfig: "regtest",
	}
	iBtcConnector, err := btc_example.NewBtcChainConnector(walletID, walletConfig)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}
	walletConnector := &bchChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   walletConfig.Currency,
			WalletType: walletConfig.Type,
		},
		ibtc:    iBtcConnector,
		chain:   &bchchaincfg.TestNet3Params,
		regtest: walletConfig.ChainConfig == "regtest",
	}

	addressesToValidate := []string{
		"bchreg:qzu7d73hakh39hc2qqmz2ln73x56hdp04cyfy5q4ye",
		"bitcoincash:qpm2qsznhks23z7629mms6s4cwef74vcwvy22gdx6a",
	}

	for tc := range addressesToValidate {
		validationResult, err := walletConnector.ValidateAddress(addressesToValidate[tc])
		fmt.Printf("validate(%s): %t", addressesToValidate[tc], validationResult)
		assert.Equal(t, true, validationResult, "Unexpected value")
		assert.Nil(t, err, "unexpected error")
	}
}

func TestBchConnector_BalanceGet(t *testing.T) {
	//Given
	walletConnector := &bchChainConnector{}

	balance, err := walletConnector.BalanceGet(models.Currency{})

	errExpected := fmt.Errorf("unsupported method: BalanceGet")

	assert.Equal(t, errExpected, err, "Unexpected error")
	assert.Empty(t, balance, "unexpected balance value")
}

func Test_nodeConnector_TxRebuild(t *testing.T) {
	walletID := uint64(1)
	walletConfig := &connector.WalletParams{
		Currency: "BCHABC",
		Type:     "hot",
		Node: connector.NodeParams{
			Host:     "104.199.25.196",
			Port:     8341,
			User:     "attic",
			Password: "8Wsujmq4JND0565itTqt",
		},
		ChainConfig: "regtest",
	}
	iBtcConnector, err := btc_example.NewBtcChainConnector(walletID, walletConfig)
	if err != nil {
		// log.Fatalf("failed to connect to Bitcoin node: %v", err.Error())
	}
	walletConnector := &bchChainConnector{
		Connector: connector.Connector{
			WalletId:   walletID,
			Currency:   walletConfig.Currency,
			WalletType: walletConfig.Type,
		},
		ibtc:    iBtcConnector,
		chain:   &bchchaincfg.TestNet3Params,
		regtest: walletConfig.ChainConfig == "regtest",
	}

	txHex := "0200000001d71f0514b1f210d374a7d5c1ea4b24bb199eb0bf1990dc9d8ec5252359b8eff600000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000e80300004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000e80300004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a00000000e803000053aeffffffff02e0f83b360000000017a914af70bbab80fb64dbf90b212f4971cc4807d0b8808700e1f505000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000"
	signatures := []connector.TxInSignatures{
		connector.TxInSignatures{
			"3044022058dbc5b8c7952fa0972d32e28d27415ede9de1c03dea74d3ae357c6f8b2c170502203558d6efdbeefb651a0be9eb5263fc5b505a842e94eabda5576022fa5f0f09c041",
			"304402201a5ff47d22d91b4c5a3195ece0c4546d49dbda9a46533bbaef885b80cf13aba002202e30e9662cbefde3c47147f7119a20e92b6d68531ffaf21344ef97fca5e8c96a41",
		},
	}
	txExpected := "0200000001d71f0514b1f210d374a7d5c1ea4b24bb199eb0bf1990dc9d8ec5252359b8eff600000000fc00473044022058dbc5b8c7952fa0972d32e28d27415ede9de1c03dea74d3ae357c6f8b2c170502203558d6efdbeefb651a0be9eb5263fc5b505a842e94eabda5576022fa5f0f09c04147304402201a5ff47d22d91b4c5a3195ece0c4546d49dbda9a46533bbaef885b80cf13aba002202e30e9662cbefde3c47147f7119a20e92b6d68531ffaf21344ef97fca5e8c96a414c695221028803d510417f3ffec81ffa81418435050d6b4693775d90a14c8abba0f74b18f42103e629b677066a100757fd930445ea5ce69d13ed4a6ee733a8e5a41f732c3311d22103fbc1d8df7237a5199dde3609ad991b12f31d9e6d09ea784eaad16c33d9f1ed6953aeffffffff02e0f83b360000000017a914af70bbab80fb64dbf90b212f4971cc4807d0b8808700e1f505000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000"

	tx, err := walletConnector.TxRebuild(txHex, signatures)
	if err != nil {
		t.Errorf("TxRebuild(): %v", err)
	}
	if tx != txExpected {
		t.Errorf("TxRebuild() = %s, want %s", tx, txExpected)
	}

	// one-time test
	//txID, err := walletConnector.TxBroadcast(tx)
	//if err != nil {
	//    t.Errorf("TxBroadcast(): %s", err.Error())
	//}
	//log.Infof("txID: %s", txID)
	// txID: dc68fa788d5a92cc7267648ad4ec36bee2e22cbcb652a753aea7b067409e3f73
}
