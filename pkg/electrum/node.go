package electrum

import (
	"net"
	"regexp"
	"strings"
)

// Node represents a node on the electrum network.
type Node struct {
	Host         string
	IP           string
	Version      string
	SSLPort      int
	TCPPort      int
	PruningLimit int
}

// NewNode constructs an instance of Node.
func NewNode(host string, IP string, version string, SSLPort int, TCPPort int, PruningLimit int) *Node {
	return &Node{Host: host, IP: IP, Version: version, SSLPort: SSLPort, TCPPort: TCPPort, PruningLimit: PruningLimit}
}

// IsValid returns true if a peer has the minimum we need to connect.
func (n *Node) IsValid() bool {
	return (ValidIP(n.IP) || ValidHostname(n.Host)) && (n.SSLPort > 0 || n.TCPPort > 0)
}

// IsOnion returns true if this node is accessible over Tor.
func (n *Node) IsOnion() bool {
	return IsOnionAddr(n.Host)
}

// SupportsTLS returns true if this host supports TLS.
func (n *Node) SupportsTLS() bool {
	return !ValidIP(n.Host) && ValidHostname(n.Host) && n.SSLPort > 0
}

// Features represents features of an electrum server.
type Features struct {
	Version      string
	SSLPort      int
	TCPPort      int
	PruningLimit int
}

// RegisterFeatures applies the listed features to the given peer.
func (n *Node) RegisterFeatures(f *Features) {
	n.Version = f.Version
	n.SSLPort = f.SSLPort
	n.TCPPort = f.TCPPort
	n.PruningLimit = f.PruningLimit
}

// IsOnionAddr checks if a hostname is a .onion address.
func IsOnionAddr(addr string) bool {
	return strings.HasSuffix(addr, ".onion")
}

// ValidIP checks if a string is an IP address that is not local / loopback.
func ValidIP(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && !(ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast())
}

// ValidHostname returns true if the hostname is valid.
func ValidHostname(hostname string) bool {
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	return re.MatchString(hostname)
}
