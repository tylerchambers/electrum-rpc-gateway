package electrum

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Peer represents a peer on the electrum network.
type Peer struct {
	Host         string
	IP           string
	Version      string
	IsOnion      bool
	SSLPort      int
	TCPPort      int
	PruningLimit int
}

// Features represents features of an electrum server.
type Features struct {
	Version      string
	SSLPort      int
	TCPPort      int
	PruningLimit int
}

// RegisterFeatures applies the listed features to the given peer.
func (p *Peer) RegisterFeatures(f *Features) {
	p.Version = f.Version
	p.SSLPort = f.SSLPort
	p.TCPPort = f.TCPPort
	p.PruningLimit = f.PruningLimit
}

// IsValid returns true if a peer has the minimum we need to connect.
func (p *Peer) IsValid() bool {
	return (ValidIP(p.IP) || ValidHostname(p.Host)) && (p.SSLPort > 0 || p.TCPPort > 0)
}

// NewPeer returns a valid electrum peer.
func NewPeer(host string, IP string, version string, isOnion bool, SSLPort int, TCPPort int, PruningLimit int) *Peer {
	return &Peer{Host: host, IP: IP, Version: version, IsOnion: isOnion, SSLPort: SSLPort, TCPPort: TCPPort, PruningLimit: PruningLimit}
}

// ServerPeersSubscriptionResp represents a response from server.peers.subscribe.
type ServerPeersSubscriptionResp struct {
	Id      string          `json:"id"`
	Version string          `json:"jsonrpc"`
	Result  [][]interface{} `json:"result"`
}

// ParseServerPeersSubscriptionResp returns a proper slice of validated peers from a ServerPeersSubscriptionResp.
func ParseServerPeersSubscriptionResp(resp *ServerPeersSubscriptionResp) ([]Peer, error) {
	if resp.Result == nil {
		return nil, errors.New("invalid message: response to parse contained a nil result field")
	}
	var peers []Peer
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
func ParsePeer(peer []interface{}) *Peer {
	newPeer := new(Peer)
	if ValidPeerResponse(peer) {
		// First element of resp. should always be an IP or an onion addr.
		if ValidIP(peer[0].(string)) {
			newPeer.IP = peer[0].(string)
		}
		// Second element is either the same IP, or a Hostname
		if !ValidIP(peer[1].(string)) && ValidHostname(peer[1].(string)) {
			newPeer.Host = peer[1].(string)
			// Set the IsOnion flag
			newPeer.IsOnion = IsOnion(newPeer.Host)
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

// ValidHostname returns true if the hostname is valid.
func ValidHostname(hostname string) bool {
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	return re.MatchString(hostname)
}

// IsOnion checks if a hostname is a .onion address.
func IsOnion(addr string) bool {
	return strings.HasSuffix(addr, ".onion")
}

// ValidIP checks if a string is an IP address that is not local / loopback.
func ValidIP(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && !(ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast())
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
