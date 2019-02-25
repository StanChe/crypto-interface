package btc

import (
	"testing"

	"github.com/stanche/crypto-interface/address/hd"
)

func TestGenerator_AddressGenerate(t *testing.T) {
	type args struct {
		params hd.GeneratorParameters
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantErr     bool
	}{
		{
			name: "empty xpubs",
			args: args{params: hd.GeneratorParameters{
				SignersXpubs: []string{},
			}},
			wantErr: true,
		},
		{
			name: "zero required",
			args: args{params: hd.GeneratorParameters{
				SignersXpubs: []string{"xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufqeGBwALkBUK", "xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ", "xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv"},
			}},
			wantErr: true,
		},
		{
			name: "empty xpub",
			args: args{params: hd.GeneratorParameters{
				SignersXpubs:    []string{"", "xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ", "xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv"},
				SignersRequired: 2,
				PathIndex:       1000,
			}},
			wantErr: true,
		},
		{
			name: "invalid xpub",
			args: args{params: hd.GeneratorParameters{
				SignersXpubs:    []string{"xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufGBwALkBUK", "xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ", "xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv"},
				SignersRequired: 2,
				PathIndex:       1000,
			}},
			wantErr: true,
		},
		{
			name: "address",
			args: args{params: hd.GeneratorParameters{
				SignersXpubs:    []string{"xpub661MyMwAqRbcEtBNvF5oTnmGFSkZvy6ShetrnbVXTz7hyKYJSNBEtKiiY9HnMeTpLKDFJRYW2QSbNGtCGdpCzwZVSPRKevufqeGBwALkBUK", "xpub661MyMwAqRbcGgsQadngKDqjvQDC299XoG8SjbpfZhKUofdVVCqehG2TCsTXNudCFyTmNL72gGmNBNbtu75Tkzz2jJMqBak8Ab71MQYs2UQ", "xpub661MyMwAqRbcFTni57UXBzWmbN3JtuoqdLivkjzkbkiPB46gDU6pYYQeE2BKRyhD1h6wXHx5jRWZh78NS45EoZPwVezgKkLjf4TTXPWh8Wv"},
				SignersRequired: 2,
				PathIndex:       1000,
			}},
			wantAddress: "2N9EsHgmGFqSUsGvBKcRqsmnWMg7dVVBYVT",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Generator{}
			gotAddress, err := g.AddressGenerate(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.AddressGenerate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("Generator.AddressGenerate() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}
