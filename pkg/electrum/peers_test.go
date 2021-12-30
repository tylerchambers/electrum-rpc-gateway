package electrum

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestIsOnion(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid onion address",
			args: args{addr: "duckduckgogg42xjoc72x3sjasowoarfbgcmvfimaftt6twagswzczad.onion"},
			want: true,
		},
		{
			name: "invalid onion address",
			args: args{addr: "duckduckgo.com"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOnion(tt.args.addr); got != tt.want {
				t.Errorf("IsOnion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPeer(t *testing.T) {
	type args struct {
		host         string
		ip           string
		version      string
		isOnion      bool
		SSLPort      int
		TCPPort      int
		PruningLimit int
	}
	tests := []struct {
		name string
		args args
		want *Peer
	}{
		{
			name: "valid new peer",
			args: args{
				host:         "electrum.blockstream.info",
				ip:           "112.123.200.12",
				version:      "v0.0.0",
				isOnion:      false,
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 0,
			},
			want: &Peer{
				Host:         "electrum.blockstream.info",
				Version:      "v0.0.0",
				IP:           "112.123.200.12",
				IsOnion:      false,
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPeer(tt.args.host, tt.args.ip, tt.args.version, tt.args.isOnion, tt.args.SSLPort, tt.args.TCPPort, tt.args.PruningLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		want *Peer
	}{
		{
			name: "parse valid peer",
			args: args{
				peer: makeRawInterfacePeer(),
			},
			want: &Peer{
				Host:         "electrum.blockstream.info",
				IP:           "232.73.129.9",
				Version:      "v0.0.1",
				IsOnion:      false,
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
		want    []Peer
		wantErr bool
	}{
		{
			name: "parse valid serverpeerssubscriptionresp",
			args: args{
				resp: validElectrumServerPeersSubscriptionRespJSON(),
			},
			want: []Peer{
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

func TestPeer_IsValid(t *testing.T) {
	type fields struct {
		Host         string
		IP           string
		Version      string
		IsOnion      bool
		SSLPort      int
		TCPPort      int
		PruningLimit int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test invalid peer",
			fields: fields{
				Host:         "",
				IP:           "",
				Version:      "",
				IsOnion:      false,
				SSLPort:      0,
				TCPPort:      0,
				PruningLimit: 0,
			},
			want: false,
		},
		{
			name: "test valid peer",
			fields: fields{
				Host:         "blockstream.electrum.info",
				IP:           "232.73.129.9",
				Version:      "v0.0.1",
				IsOnion:      false,
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 0,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Peer{
				Host:         tt.fields.Host,
				IP:           tt.fields.IP,
				Version:      tt.fields.Version,
				IsOnion:      tt.fields.IsOnion,
				SSLPort:      tt.fields.SSLPort,
				TCPPort:      tt.fields.TCPPort,
				PruningLimit: tt.fields.PruningLimit,
			}
			if got := p.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeer_RegisterFeatures(t *testing.T) {
	type fields struct {
		Host         string
		IP           string
		Version      string
		IsOnion      bool
		SSLPort      int
		TCPPort      int
		PruningLimit int
	}
	type args struct {
		f *Features
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Peer
	}{
		{
			name: "test valid feature registration",
			fields: fields{
				Host:         "electrum.blockstream.info",
				IP:           "232.73.129.9",
				Version:      "",
				IsOnion:      false,
				SSLPort:      0,
				TCPPort:      0,
				PruningLimit: 0,
			},
			args: args{
				f: &Features{
					Version:      "v0.0.1",
					SSLPort:      50002,
					TCPPort:      50001,
					PruningLimit: 5,
				},
			},
			want: &Peer{
				Host:         "electrum.blockstream.info",
				IP:           "232.73.129.9",
				Version:      "v0.0.1",
				IsOnion:      false,
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Peer{
				Host:         tt.fields.Host,
				IP:           tt.fields.IP,
				Version:      tt.fields.Version,
				IsOnion:      tt.fields.IsOnion,
				SSLPort:      tt.fields.SSLPort,
				TCPPort:      tt.fields.TCPPort,
				PruningLimit: tt.fields.PruningLimit,
			}
			p.RegisterFeatures(tt.args.f)
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Peer.RegisterFeatures() got = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestValidHostname(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid hostname",
			args: args{
				hostname: "google.com",
			},
			want: true,
		},
		{
			name: "invalid hostname",
			args: args{
				hostname: "👍___+",
			},
			want: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidHostname(tt.args.hostname); got != tt.want {
				t.Errorf("ValidHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidIP(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid IP",
			args: args{
				s: "232.73.129.9",
			},
			want: true,
		},
		{
			name: "valid loopback ip",
			args: args{
				s: "127.0.0.1",
			},
			want: false,
		},
		{
			name: "valid internal ip",
			args: args{
				s: "192.168.0.1",
			},
			want: false,
		},
		{
			name: "not an IP",
			args: args{
				s: "hello, world!",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidIP(tt.args.s); got != tt.want {
				t.Errorf("ValidIP() = %v, want %v", got, tt.want)
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
