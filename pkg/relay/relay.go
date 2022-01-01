package relay

import (
	"github.com/tylerchambers/electrumrelay/pkg/electrum"
	"log"
	"time"
)

// Relay represents an electrum relay.
// Handles the logic of finding peers, then taking and forwarding requests to them.
type Relay struct {
	Peers            []electrum.Node
	ForbiddenMethods []string
}

// Init takes an initial peer, and gets info of other peers.
func (r *Relay) Init(initialPeer *electrum.Node) error {
	c := electrum.Client{
		InfoLogger:    log.Default(),
		WarningLogger: log.Default(),
		ErrorLogger:   log.Default(),
	}
	peers, err := c.GetPeerInfo(initialPeer, 444, time.Second*10)
	if err != nil {
		return err
	}
	r.Peers = peers
	return nil
}
