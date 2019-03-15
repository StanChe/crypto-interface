package signers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestBtcSigner_Sign(t *testing.T) {
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

				signer := NewBtcSigner("BTC", nil, nil, nil)
				codeExpected := "BTC"

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
				signer := NewBtcSigner("BTC", kp, &BtcNetParams, nil)

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
			"Positive signing tx with one input (signer-a)",
			func(t *testing.T) {
				kp := New(component1)
				signer := NewBtcSigner("BTC", kp, &BtcNetParams, BtcTxInputSignature)

				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")

				txSign, _ := hex.DecodeString("30440220596c276e66186b98e1b190a626a94b30760718c99f2db32d2e165e7075c3f67302207fd6cd72995239952769b7ff2c61f4e952a05a3a9970f564f09f3efe51feded201")
				signStruct := btcSignature{
					Ind: 0,
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
			"Positive signing tx with one input (signer-b)",
			func(t *testing.T) {
				kp := New(component2)
				signer := NewBtcSigner("BTC", kp, &BtcNetParams, BtcTxInputSignature)

				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")

				txSign, _ := hex.DecodeString("3045022100cc08a8be0f1021f9029b0fd428a0d1575c39e215ee396672eb70dd350f5b17d30220075c8eaf70f6d0dfcdf379c4a7f98ae7eb64ef765ac8d7be1b30f6e7e4c4181301")
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
				signer := NewBtcSigner("BTC", kp, &BtcNetParams, BtcTxInputSignature)

				txData := []byte("0200000001db172762bebe28c7f79bcea59647ca37e4e38603618bebbf8407bf44b727c58f00000000fd16010001ff01ff4d0e01524c57ff0488b21e0000000000000000002231c2b6a33377bc6fb0806268e3627602987340ed2c5e6be0d7be7f24161bae038b8001ff63faf92876effaa8cb774ee8a7260b014922607e191b22fb88d3ef1700000000020000004c57ff0488b21e000000000000000000d77de533cea4f03402d513aa6b682cd1a69409564a6c4cddb37c8eed4705d0c603d2a614051301da597eea74316d7e404d89d5eb850238c2c1b3d536c5d5c07a5900000000020000004c57ff0488b21e0000000000000000005c65a74ec6c4922e3df98f50f7c297f62477d123989d9c69ad7de1322cc8394c02cc24a901a51e4e1525343049f11ded77391bf579bc020f08e6956a6eadb13b5a000000000200000053aeffffffff019c35f8030000000017a9140a4aa12d8ff4bf38647a21bb9f72c3602fecaa448700000000")

				txSign, _ := hex.DecodeString("304402204c98e1508cd33482d004ec5044f71128c0adf5fbcf2fd3da0976023a04811ba20220039926ee4a5449a3ab44fdae71a356eed7dcf5d37c28bf79bfa11f67f4ee444c01")
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
