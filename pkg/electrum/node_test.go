package electrum

import (
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
			if got := IsOnionAddr(tt.args.addr); got != tt.want {
				t.Errorf("IsOnionAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	type args struct {
		host         string
		ip           string
		version      string
		SSLPort      int
		TCPPort      int
		PruningLimit int
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "valid new peer",
			args: args{
				host:         "electrum.blockstream.info",
				ip:           "112.123.200.12",
				version:      "v0.0.0",
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 0,
			},
			want: &Node{
				Host:         "electrum.blockstream.info",
				Version:      "v0.0.0",
				IP:           "112.123.200.12",
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.host, tt.args.ip, tt.args.version, tt.args.SSLPort, tt.args.TCPPort, tt.args.PruningLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_IsValid(t *testing.T) {
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
			p := &Node{
				Host:         tt.fields.Host,
				IP:           tt.fields.IP,
				Version:      tt.fields.Version,
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

func TestNode_RegisterFeatures(t *testing.T) {
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
		want   *Node
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
			want: &Node{
				Host:         "electrum.blockstream.info",
				IP:           "232.73.129.9",
				Version:      "v0.0.1",
				SSLPort:      50002,
				TCPPort:      50001,
				PruningLimit: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Node{
				Host:         tt.fields.Host,
				IP:           tt.fields.IP,
				Version:      tt.fields.Version,
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
