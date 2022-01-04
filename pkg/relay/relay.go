package relay

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tylerchambers/electrumrelay/pkg/electrum"
)

// Relay represents an electrum relay.
// Handles the logic of finding peers, then taking and forwarding requests to them.
type Relay struct {
	Peers            []electrum.Node
	PeerMutex        sync.Mutex
	ForbiddenMethods []string
	ElectrumClient   *electrum.Client
}

// NewRelay constructs a new JSON RPC Relay.
func NewRelay(peers []electrum.Node, forbiddenMethods []string, electrumClient *electrum.Client) *Relay {
	return &Relay{Peers: peers, PeerMutex: sync.Mutex{}, ForbiddenMethods: forbiddenMethods, ElectrumClient: electrumClient}
}

func (r *Relay) NoOnions() []electrum.Node {
	var out []electrum.Node
	for _, v := range r.Peers {
		if !v.IsOnion() {
			out = append(out, v)
		}
	}
	return out
}

// Bootstrap takes an initial peer, asks for its peers, then registers them.
func (r *Relay) Bootstrap(initialPeer *electrum.Node) error {
	// Make a random request ID
	peers, err := r.ElectrumClient.GetPeerInfo(initialPeer, rand.Intn(512), time.Second*10)
	if err != nil {
		return err
	}
	err = r.RegisterPeers(peers)
	if err != nil {
		return err
	}
	return nil
}

// RegisterPeer adds a peer to the relay's slice of peers.
func (r *Relay) RegisterPeer(peer *electrum.Node) error {
	if !peer.IsValid() {
		return fmt.Errorf("invalid peer: %s not registering", peer.Host)
	}
	r.PeerMutex.Lock()
	r.Peers = append(r.Peers, *peer)
	r.PeerMutex.Unlock()
	return nil
}

// RegisterPeers adds a slice of peers to the relay's slice of peers.
func (r *Relay) RegisterPeers(peers []electrum.Node) error {
	for _, peer := range peers {
		if !peer.IsValid() {
			return fmt.Errorf("invalid peer: %s not registering any", peer.Host)
		}
	}
	r.PeerMutex.Lock()
	r.Peers = append(r.Peers, peers...)
	r.PeerMutex.Unlock()
	return nil
}

// ValidateRequest validates an incoming HTTP Request, parses out the JSON RPC request it contains in the body, and
// makes sure it's allowed.
func (r *Relay) ValidateRequest(req *http.Request) ([]byte, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse electrum request: %v", err)
	}
	if !r.AllowedMethod(b) {
		return nil, fmt.Errorf("method is forbidden by the relay")
	}
	return b, nil
}

// AllowedMethod returns true if the method it is passed is not forbidden by the relay.
func (r *Relay) AllowedMethod(req []byte) bool {
	s := string(req)
	for _, v := range r.ForbiddenMethods {
		if strings.Contains(v, s) {
			return false
		}
	}
	return true
}

// RandomNode selects and returns a random electrum node from the list of peers.
// holdTheOnions only returns clearnet nodes until Tor support is implemented.
func (r *Relay) RandomNode(holdTheOnions bool) *electrum.Node {
	if holdTheOnions {
		r.PeerMutex.Lock()
		noOnions := r.NoOnions()
		random := &noOnions[rand.Intn(len(noOnions))]
		r.PeerMutex.Unlock()
		return random
	}
	r.PeerMutex.Lock()
	random := &r.Peers[rand.Intn(len(r.Peers))]
	r.PeerMutex.Unlock()
	return random
}

// ForwardRequest forwards the request to a random peer, and returns the response as bytes.
func (r *Relay) ForwardRequest(req []byte) ([]byte, error) {
	n := r.RandomNode(true)
	resp, err := r.ElectrumClient.SendRequestBytes(req, n, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("error forwarding request %s to node %s %v", string(req), n.Host, err)
	}
	return resp, nil
}
