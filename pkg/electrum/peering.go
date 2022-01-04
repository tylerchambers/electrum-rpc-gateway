package electrum

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// NewPeerRequest is a convenience function for creation a server.peers.subscribe request.
func NewPeerRequest(id int) *JSONRPCRequest {
	req := new(JSONRPCRequest)
	reqStr := fmt.Sprintf(`{"jsonrpc":"2.0","method":"server.peers.subscribe","params":[],"id":%d}`, id)
	_ = json.Unmarshal([]byte(reqStr), req)
	return req
}

// ServerPeersSubscriptionResp represents a response from server.peers.subscribe.
type ServerPeersSubscriptionResp struct {
	ID      int             `json:"id"`
	Version string          `json:"jsonrpc"`
	Result  [][]interface{} `json:"result"`
}

// ParseServerPeersSubscriptionResp returns a proper slice of validated peers from a ServerPeersSubscriptionResp.
func ParseServerPeersSubscriptionResp(resp *ServerPeersSubscriptionResp) ([]Node, error) {
	if resp.Result == nil {
		return nil, errors.New("invalid message: response to parse contained a nil result field")
	}
	var peers []Node
	for _, peer := range resp.Result {
		p := ParsePeer(peer)
		if p != nil {
			peers = append(peers, *p)
		}
	}
	if len(peers) == 0 {
		return nil, errors.New("unable to parse valid peers from the response")
	}
	return peers, nil
}

// ParsePeer parses a peer from the array of peers in the response from a server.peers.subscribe request.
// Returns nil if the peer is invalid.
func ParsePeer(peer []interface{}) *Node {
	newPeer := new(Node)
	if ValidPeerResponse(peer) {
		// First element of resp. should always be an IP or an onion addr.
		if ValidIP(peer[0].(string)) {
			newPeer.IP = peer[0].(string)
		}
		// Second element is either the same IP, or a Hostname
		if !ValidIP(peer[1].(string)) && ValidHostname(peer[1].(string)) {
			newPeer.Host = peer[1].(string)
		}
		// Third element is another array of info.
		features := ParseServerFeatures(peer[2].([]interface{}))
		newPeer.RegisterFeatures(features)
	}
	if newPeer.IsValid() {
		return newPeer
	}
	return nil
}

// ValidPeerResponse checks if the array representing the peer information contains any nil fields.
func ValidPeerResponse(peer []interface{}) bool {
	// Proper response is 3 long as defined here:
	// https://electrumx-spesmilo.readthedocs.io/en/latest/protocol-methods.html#server-peers-subscribe
	if len(peer) != 3 {
		return false
	}
	// No nil values allowed
	for _, v := range peer {
		if v == nil {
			return false
		}
	}
	return true
}

// ParseServerFeatures parses the third element of the response array, containing server features.
func ParseServerFeatures(r []interface{}) *Features {
	f := new(Features)

	for _, v := range r {
		// These should all be non-unicode, so this should work.
		switch rune(v.(string)[0]) {
		case 'v':
			f.Version = v.(string)
		case 'p':
			limit, err := strconv.Atoi(v.(string)[1:])
			if err != nil && limit <= 0 {
				break
			}
			f.PruningLimit = limit
		case 't':
			port, err := strconv.Atoi(v.(string)[1:])
			if err != nil || port <= 0 {
				break
			}
			f.TCPPort = port
		case 's':
			port, err := strconv.Atoi(v.(string)[1:])
			if err != nil || port <= 0 {
				continue
			}
			f.SSLPort = port
		}
	}
	return f
}
