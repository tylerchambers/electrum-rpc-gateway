package relay

import (
	"reflect"
	"sync"
	"testing"

	"github.com/tylerchambers/electrumrelay/pkg/electrum"
)

func TestNewRelay(t *testing.T) {
	type args struct {
		peers            []electrum.Node
		forbiddenMethods []string
		electrumClient   *electrum.Client
	}
	tests := []struct {
		name string
		args args
		want *Relay
	}{
		{
			name: "Test NewRelay",
			args: args{
				peers: []electrum.Node{
					{
						Host:         "electrum.blockstream.info",
						IP:           "",
						Version:      "",
						SSLPort:      50002,
						TCPPort:      0,
						PruningLimit: 0,
					},
				},
				forbiddenMethods: []string{
					"blockchain.scripthash.subscribe",
				},
				electrumClient: &electrum.Client{
					InfoLogger:    nil,
					WarningLogger: nil,
					ErrorLogger:   nil,
				},
			},
			want: &Relay{
				Peers: []electrum.Node{
					{
						Host:         "electrum.blockstream.info",
						IP:           "",
						Version:      "",
						SSLPort:      50002,
						TCPPort:      0,
						PruningLimit: 0,
					},
				},
				PeerMutex: sync.Mutex{},
				ForbiddenMethods: []string{
					"blockchain.scripthash.subscribe",
				},
				ElectrumClient: &electrum.Client{
					InfoLogger:    nil,
					WarningLogger: nil,
					ErrorLogger:   nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRelay(tt.args.peers, tt.args.forbiddenMethods, tt.args.electrumClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelay_NoOnions(t *testing.T) {
	type fields struct {
		Peers            []electrum.Node
		PeerMutex        sync.Mutex
		ForbiddenMethods []string
		ElectrumClient   *electrum.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []electrum.Node
	}{
		{
			name: "Test NoOnions",
			fields: fields{
				Peers: []electrum.Node{
					{
						Host: "electrum.blockstream.info",
					},
					{
						Host: "wsw6tua3xl24gsmi264zaep6seppjyrkyucpsmuxnjzyt3f3j6swshad.onion",
					},
				},
			},
			want: []electrum.Node{
				{
					Host: "electrum.blockstream.info",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relay{
				Peers:            tt.fields.Peers,
				PeerMutex:        tt.fields.PeerMutex,
				ForbiddenMethods: tt.fields.ForbiddenMethods,
				ElectrumClient:   tt.fields.ElectrumClient,
			}
			if got := r.NoOnions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Relay.NoOnions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelay_Bootstrap(t *testing.T) {
	type fields struct {
		Peers            []electrum.Node
		PeerMutex        sync.Mutex
		ForbiddenMethods []string
		ElectrumClient   *electrum.Client
	}
	type args struct {
		initialPeer *electrum.Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relay{
				Peers:            tt.fields.Peers,
				PeerMutex:        tt.fields.PeerMutex,
				ForbiddenMethods: tt.fields.ForbiddenMethods,
				ElectrumClient:   tt.fields.ElectrumClient,
			}
			if err := r.Bootstrap(tt.args.initialPeer); (err != nil) != tt.wantErr {
				t.Errorf("Relay.Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelay_RegisterPeer(t *testing.T) {
	type fields struct {
		Peers            []electrum.Node
		PeerMutex        sync.Mutex
		ForbiddenMethods []string
		ElectrumClient   *electrum.Client
	}
	type args struct {
		peer *electrum.Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relay{
				Peers:            tt.fields.Peers,
				PeerMutex:        tt.fields.PeerMutex,
				ForbiddenMethods: tt.fields.ForbiddenMethods,
				ElectrumClient:   tt.fields.ElectrumClient,
			}
			if err := r.RegisterPeer(tt.args.peer); (err != nil) != tt.wantErr {
				t.Errorf("Relay.RegisterPeer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelay_RegisterPeers(t *testing.T) {
	type fields struct {
		Peers            []electrum.Node
		PeerMutex        sync.Mutex
		ForbiddenMethods []string
		ElectrumClient   *electrum.Client
	}
	type args struct {
		peers []electrum.Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relay{
				Peers:            tt.fields.Peers,
				PeerMutex:        tt.fields.PeerMutex,
				ForbiddenMethods: tt.fields.ForbiddenMethods,
				ElectrumClient:   tt.fields.ElectrumClient,
			}
			if err := r.RegisterPeers(tt.args.peers); (err != nil) != tt.wantErr {
				t.Errorf("Relay.RegisterPeers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelay_AllowedMethod(t *testing.T) {
	type fields struct {
		Peers            []electrum.Node
		PeerMutex        sync.Mutex
		ForbiddenMethods []string
		ElectrumClient   *electrum.Client
	}
	type args struct {
		req []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "forbidden method detected in req",
			fields: fields{
				ForbiddenMethods: []string{"server.version"},
			},
			args: args{
				req: []byte("server.version"),
			},
			want: false,
		},
		{
			name: "forbidden method detected in req",
			fields: fields{
				ForbiddenMethods: []string{"server.version"},
			},
			args: args{
				req: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque feugiat pulvinar urna, sit amet luctus mi mattis at"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relay{
				Peers:            tt.fields.Peers,
				PeerMutex:        tt.fields.PeerMutex,
				ForbiddenMethods: tt.fields.ForbiddenMethods,
				ElectrumClient:   tt.fields.ElectrumClient,
			}
			if got := r.AllowedMethod(tt.args.req); got != tt.want {
				t.Errorf("Relay.AllowedMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
