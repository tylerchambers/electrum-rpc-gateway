package electrum

import (
	"encoding/json"
	"reflect"
	"testing"
)

// This is necessary for some tests.
func makeRawInterfacePeer() []interface{} {
	out := make([]interface{}, 3)
	out[0] = "232.73.129.9"
	out[1] = "electrum.blockstream.info"
	features := make([]interface{}, 4)
	features[0] = "v0.0.1"
	features[1] = "p10000"
	features[2] = "t50002"
	features[3] = "s50001"
	out[2] = features
	return out
}

func TestParsePeer(t *testing.T) {
	type args struct {
		peer []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "parse valid peer",
			args: args{
				peer: makeRawInterfacePeer(),
			},
			want: &Node{
				Host:         "electrum.blockstream.info",
				IP:           "232.73.129.9",
				Version:      "v0.0.1",
				SSLPort:      50001,
				TCPPort:      50002,
				PruningLimit: 10000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePeer(tt.args.peer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeRawInterfaceFeatures() []interface{} {
	features := make([]interface{}, 4)
	features[0] = "v0.0.1"
	features[1] = "p10000"
	features[2] = "t50002"
	features[3] = "s50001"
	return features
}

func TestParseServerFeatures(t *testing.T) {
	type args struct {
		r []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Features
	}{
		{
			name: "parse valid features",
			args: args{
				r: makeRawInterfaceFeatures(),
			},
			want: &Features{
				Version:      "v0.0.1",
				SSLPort:      50001,
				TCPPort:      50002,
				PruningLimit: 10000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseServerFeatures(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseServerFeatures() = %v, want %v", got, tt.want)
			}
		})
	}
}

func validElectrumServerPeersSubscriptionRespJSON() *ServerPeersSubscriptionResp {
	str := `{"id":"","jsonrpc":"2.0","result":[["164.132.182.11","1.fulcrum-node.com",["v1.4.5","s50002"]],["213.152.106.56","electrum.pabu.io",["v1.4.2","s50002"]],["144.76.84.234","electrum.jochen-hoenicke.de",["v1.4.5","s50006","t50099"]]]}`
	out := ServerPeersSubscriptionResp{}
	err := json.Unmarshal([]byte(str), &out)
	if err != nil {
		panic(err)
	}
	return &out
}

func TestParseServerPeersSubscriptionResp(t *testing.T) {
	type args struct {
		resp *ServerPeersSubscriptionResp
	}
	tests := []struct {
		name    string
		args    args
		want    []Node
		wantErr bool
	}{
		{
			name: "parse valid serverpeerssubscriptionresp",
			args: args{
				resp: validElectrumServerPeersSubscriptionRespJSON(),
			},
			want: []Node{
				{
					Host:    "1.fulcrum-node.com",
					IP:      "164.132.182.11",
					Version: "v1.4.5",
					SSLPort: 50002,
				},
				{
					Host:    "electrum.pabu.io",
					IP:      "213.152.106.56",
					Version: "v1.4.2",
					SSLPort: 50002,
				},
				{
					Host:    "electrum.jochen-hoenicke.de",
					IP:      "144.76.84.234",
					Version: "v1.4.5",
					SSLPort: 50006,
					TCPPort: 50099,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseServerPeersSubscriptionResp(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseServerPeersSubscriptionResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseServerPeersSubscriptionResp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidPeerResponse(t *testing.T) {
	type args struct {
		peer []interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid peer response interface",
			args: args{
				peer: makeRawInterfacePeer(),
			},
			want: true,
		},
		{
			name: "valid peer response interface",
			args: args{
				peer: []interface{}{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidPeerResponse(tt.args.peer); got != tt.want {
				t.Errorf("ValidPeerResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
