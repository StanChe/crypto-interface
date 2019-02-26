package btc_example

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/stanche/crypto-interface/connector"
	"github.com/stanche/crypto-interface/importer"
)

func TestNewBlockChainImporter(t *testing.T) {
	type args struct {
		node        importer.NodeParams
		chainParams chaincfg.Params
		txBatchSize int
	}
	tests := []struct {
		name    string
		args    args
		want    importer.BlockChainImporter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBlockChainImporter(tt.args.node, tt.args.chainParams, tt.args.txBatchSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBlockChainImporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlockChainImporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBtcBlockChainImporter_getBlockByNumber(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		number uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantBlock *wire.MsgBlock
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			gotBlock, err := bci.getBlockByNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("BtcBlockChainImporter.getBlockByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBlock, tt.wantBlock) {
				t.Errorf("BtcBlockChainImporter.getBlockByNumber() = %v, want %v", gotBlock, tt.wantBlock)
			}
		})
	}
}

func TestBtcBlockChainImporter_GetBlockHashesByNumber(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		number uint64
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantHash     string
		wantPrevHash string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			gotHash, gotPrevHash, err := bci.GetBlockHashesByNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("BtcBlockChainImporter.GetBlockHashesByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("BtcBlockChainImporter.GetBlockHashesByNumber() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
			if gotPrevHash != tt.wantPrevHash {
				t.Errorf("BtcBlockChainImporter.GetBlockHashesByNumber() gotPrevHash = %v, want %v", gotPrevHash, tt.wantPrevHash)
			}
		})
	}
}

func TestBtcBlockChainImporter_ProcessBlock(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		blockNumber uint64
		currencies  []connector.Currency
		addresses   []connector.Address
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantOperations []importer.Operation
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			gotOperations, err := bci.ProcessBlock(tt.args.blockNumber, tt.args.currencies, tt.args.addresses)
			if (err != nil) != tt.wantErr {
				t.Errorf("BtcBlockChainImporter.ProcessBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOperations, tt.wantOperations) {
				t.Errorf("BtcBlockChainImporter.ProcessBlock() = %v, want %v", gotOperations, tt.wantOperations)
			}
		})
	}
}

func TestBtcBlockChainImporter_min(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			if got := bci.min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("BtcBlockChainImporter.min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBtcBlockChainImporter_isAddressInList(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		target    string
		addresses []connector.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			if got := bci.isAddressInList(tt.args.target, tt.args.addresses); got != tt.want {
				t.Errorf("BtcBlockChainImporter.isAddressInList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBtcBlockChainImporter_processTransaction(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		d processTxData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   processTxResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			if got := bci.processTransaction(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BtcBlockChainImporter.processTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBtcBlockChainImporter_parseOutputs(t *testing.T) {
	type fields struct {
		client      *rpcclient.Client
		chainParams chaincfg.Params
		txBatchSize int
	}
	type args struct {
		txOuts []*wire.TxOut
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []outputParsed
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bci := BtcBlockChainImporter{
				client:      tt.fields.client,
				chainParams: tt.fields.chainParams,
				txBatchSize: tt.fields.txBatchSize,
			}
			got, err := bci.parseOutputs(tt.args.txOuts)
			if (err != nil) != tt.wantErr {
				t.Errorf("BtcBlockChainImporter.parseOutputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BtcBlockChainImporter.parseOutputs() = %v, want %v", got, tt.want)
			}
		})
	}
}
