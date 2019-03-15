package signers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestBchSigner_Sign(t *testing.T) {
	type fields struct {
		Net     *chaincfg.Params
		KeyData []byte
	}
	type args struct {
		txData     []byte
		signParams *[]uint64
	}

	// #1
	component1, _ := hex.DecodeString("0635671834e54c61b9352f26595d9615ef1e5840c7f64af198e4a10ed7140dd0")

	// #2
	component2, _ := hex.DecodeString("b918edc07dd94ad9b8f705cddc6d133bfbe3aa9bdaca4c1fb99c755ff222d461")

	// #3
	component3, _ := hex.DecodeString("1c4798b1fa6841e4b2c034c77d9221bdf44b0738f47149d88b40f772866c3649")

	cases := []struct {
		name string
		exec func(t *testing.T)
	}{
		{
			"Positive getting CurrencyType",
			func(t *testing.T) {

				signer := NewBtcSigner("BCH", nil, nil, nil)
				codeExpected := "BCH"

				//When
				code := signer.CurrencyType()

				//Then
				assert.Equal(t, codeExpected, code, "unexpected CyrrencyType")
			},
		},
		{
			"Positive getting XPub",
			func(t *testing.T) {
				kp := New(component1)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, nil)

				xpubExpected := BtcPublicAttributes{
					XPub: "xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufqeGBwALkBUK",
				}

				//When
				xpub, err := signer.Public()

				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, xpubExpected, xpub, "unexpected status")
			},
		},
		{
			"Positive integration test (signer-a)",
			func(t *testing.T) {
				kp := New(component1)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, BchTxInputSignature)

				txData := []byte("0200000001d71f0514b1f210d374a7d5c1ea4b24bb199eb0bf1990dc9d8ec5252359b8eff600000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000e80300004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000e80300004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a00000000e803000053aeffffffff02e0f83b360000000017a914af70bbab80fb64dbf90b212f4971cc4807d0b8808700e1f505000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000")
				signParams := []uint64{
					1010000000, // 10.1 BCH
				}

				txSign, _ := hex.DecodeString("3044022058dbc5b8c7952fa0972d32e28d27415ede9de1c03dea74d3ae357c6f8b2c170502203558d6efdbeefb651a0be9eb5263fc5b505a842e94eabda5576022fa5f0f09c041")
				signStruct := btcSignature{
					Ind: 0,
					Val: txSign,
				}
				byteSignature, _ := json.Marshal(signStruct)
				signature := base64.StdEncoding.EncodeToString(byteSignature)

				signExpected := []string{
					signature,
				}
				sign, err := signer.Sign(txData, signParams)
				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, signExpected, sign, "unexpected signature")
			},
		},
		{
			"Positive integration test (signer-b)",
			func(t *testing.T) {
				kp := New(component2)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, BchTxInputSignature)

				txData := []byte("0200000001d71f0514b1f210d374a7d5c1ea4b24bb199eb0bf1990dc9d8ec5252359b8eff600000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000e80300004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000e80300004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a00000000e803000053aeffffffff02e0f83b360000000017a914af70bbab80fb64dbf90b212f4971cc4807d0b8808700e1f505000000001976a914b9e6fa37edaf12df0a0036257e7e89a9abb42fae88ac00000000")
				signParams := []uint64{
					1010000000, // 10.1 BCH
				}

				txSign, _ := hex.DecodeString("304402201a5ff47d22d91b4c5a3195ece0c4546d49dbda9a46533bbaef885b80cf13aba002202e30e9662cbefde3c47147f7119a20e92b6d68531ffaf21344ef97fca5e8c96a41")
				signStruct := btcSignature{
					Ind: 1,
					Val: txSign,
				}
				byteSignature, _ := json.Marshal(signStruct)
				signature := base64.StdEncoding.EncodeToString(byteSignature)

				signExpected := []string{
					signature,
				}
				sign, err := signer.Sign(txData, signParams)
				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, signExpected, sign, "unexpected signature")
			},
		},
		{
			"Positive signing tx with one input (signer-a)",
			func(t *testing.T) {
				kp := New(component1)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, BchTxInputSignature)
				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")
				signParams := []uint64{
					110000000, // 1.1 BCH
				}

				txSign, _ := hex.DecodeString("30450221008fa912b8adb46e09eace525c90d3050f6207ed706eb49558d44e585f89c6d4d6022067c18ce15408a4d1d165093be0c9f376cc2c48b623a00548c662c131bed0020641")
				signStruct := btcSignature{
					Ind: 0,
					Val: txSign,
				}
				byteSignature, _ := json.Marshal(signStruct)
				signature := base64.StdEncoding.EncodeToString(byteSignature)

				signExpected := []string{
					signature,
				}
				sign, err := signer.Sign(txData, signParams)
				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, signExpected, sign, "unexpected signature")
			},
		},
		{
			"Positive signing tx with one input (signer-b)",
			func(t *testing.T) {
				kp := New(component2)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, BchTxInputSignature)
				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")

				txSign, _ := hex.DecodeString("3045022100c446a6f6281548c2bd11906b9c53d8ad88c1f7ad6124b1ed81e5b35a2fb6efc2022079a9dccd8d19c3abc2cf9bca91ab3fc23552f2a37c5fea850d771e734e475fcc41")
				signStruct := btcSignature{
					Ind: 2,
					Val: txSign,
				}
				byteSignature, _ := json.Marshal(signStruct)
				signature := base64.StdEncoding.EncodeToString(byteSignature)

				signExpected := []string{
					signature,
				}
				sign, err := signer.Sign(txData, nil)
				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, signExpected, sign, "unexpected signature")
			},
		},
		{
			"Positive signing tx with one input (signer-c)",
			func(t *testing.T) {
				kp := New(component3)
				signer := NewBtcSigner("BCH", kp, &BtcNetParams, BchTxInputSignature)

				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")

				txSign, _ := hex.DecodeString("304402205be45f0d347e9d454a529a01d3493675b9d8be844908c38d43b3647ecaad412e022037d8e923fa8c12c374e5530bf2b4b29ddd7a355430c8d543b9884b04543d986841")
				signStruct := btcSignature{
					Ind: 1,
					Val: txSign,
				}
				byteSignature, _ := json.Marshal(signStruct)
				signature := base64.StdEncoding.EncodeToString(byteSignature)

				signExpected := []string{
					signature,
				}
				sign, err := signer.Sign(txData, nil)
				//Then
				assert.Equal(t, nil, err, "unexpected error")
				assert.Equal(t, signExpected, sign, "unexpected signature")
			},
		},
		//
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
