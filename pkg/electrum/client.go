package electrum

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// Client handles connecting, sending requests, and the logging thereof to Electrum servers.
type Client struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

// Connect tries to connect to a node in the following order: Tor, TLS, TCP.
func (c *Client) Connect(n *Node, timeout time.Duration) (net.Conn, error) {
	if n.IsOnion() {
		c.ErrorLogger.Printf("failed to connect to %s: tor support not yet implemented\n", n.Host)
		return nil, errors.New("to support not yet implemented")
	}
	if n.SupportsTLS() {
		c.InfoLogger.Printf("%s supports TLS, attempting TLS connection\n", n.Host)
		conn, err := c.GetTLSConn(n, timeout)
		if err != nil {
			c.ErrorLogger.Printf("error establishing TLS connection to: %s\n", n.Host)
			connErr := conn.Close()
			if connErr != nil {
				c.ErrorLogger.Printf("could not close connection to: %s after failed TLS connection attempt: %v\n", n.Host, connErr)
			}
			return nil, err
		}
		return conn, nil
	}
	conn, err := c.GetConn(n, timeout)
	c.InfoLogger.Printf("%s supports TCP, attempting TCP connection\n", n.Host)
	if err != nil {
		c.ErrorLogger.Printf("error establishing TLS connection to: %s\n: %v", n.Host, err)
		connErr := conn.Close()
		if connErr != nil {
			c.ErrorLogger.Printf("could not close connection to: %s after failed TCP connection attempt: %v\n", n.Host, connErr)
		}
		return nil, err
	}
	return conn, nil
}

// GetTLSConn establishes a TLS connection to a given node.
func (c *Client) GetTLSConn(n *Node, timeout time.Duration) (*tls.Conn, error) {
	if n.IsOnion() {
		c.ErrorLogger.Printf("failed to connect to %s: tor support not yet implemented\n", n.Host)
		return nil, errors.New("tor support not yet implemented")
	}
	if !n.SupportsTLS() {
		c.ErrorLogger.Printf("%s does not support TLS, not attempting to connect\n", n.Host)
		return nil, errors.New("node does not support SSL/TLS")
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	connStr := fmt.Sprintf("%s:%d", n.Host, n.SSLPort)
	conn, err := tls.DialWithDialer(dialer, "tcp", connStr, conf)
	if err != nil {
		c.ErrorLogger.Printf("error establishing TLS connection to: %s\n: %v", n.Host, err)
		return nil, fmt.Errorf("could not establish TLS connection to %s: %v", connStr, err)
	}
	c.InfoLogger.Printf("successfully established TLS connection to %s\n", connStr)
	return conn, nil
}

// GetConn establishes a TCP connection to a given node.
func (c *Client) GetConn(n *Node, timeout time.Duration) (net.Conn, error) {
	if n.IsOnion() {
		c.ErrorLogger.Printf("failed to connect to %s: tor support not yet implemented\n", n.Host)
		return nil, errors.New("tor support not yet implemented")
	}
	connStr := fmt.Sprintf("%s:%d", n.Host, n.TCPPort)
	c.InfoLogger.Printf("establishing TCP connection to %s\n", connStr)
	conn, err := net.DialTimeout("tcp", connStr, timeout)
	if err != nil {
		c.ErrorLogger.Printf("could not establish TLS connection to %s: %v\n", connStr, err)
		return nil, fmt.Errorf("could not establish TLS connection to %s: %v", connStr, err)
	}
	c.InfoLogger.Printf("successfully established TCP connection to %s\n", connStr)
	return conn, nil
}

// SendRequest sends a JSON RPC Request to a node, and returns a response as bytes.
func (c *Client) SendRequest(req *JSONRPCRequest, n *Node, timeout time.Duration) ([]byte, error) {
	c.InfoLogger.Printf("attempting to connect to %s\n", n.Host)
	conn, err := c.Connect(n, timeout)
	c.InfoLogger.Printf("sending request ID: %s to: %s\n", req.ID, n.Host)
	resp, err := req.Send(conn)
	if err != nil {
		c.ErrorLogger.Printf("error sending request ID: %s to: %s: %v\n", req.ID, n.Host, err)
		connErr := conn.Close()
		if connErr != nil {
			c.ErrorLogger.Printf("could not close connection to: %s after failed request ID: %s: %v\n", n.Host, req.ID, connErr)
		}
		return nil, err
	}
	_ = conn.Close()
	return resp, nil
}

// GetPeerInfo gets peer information from a node by sending it a server.peers.subscribe JSON RPC Request
// It then parses the response and returns a []Node of Electrum peers.
func (c *Client) GetPeerInfo(n *Node, reqID int, timeout time.Duration) ([]Node, error) {
	if n.IsOnion() {
		c.ErrorLogger.Printf("failed to connect to %s: tor support not yet implemented\n", n.Host)
		return nil, errors.New("tor support not yet implemented")
	}
	resp, err := c.SendRequest(NewPeerRequest(reqID), n, timeout)
	if err != nil {
		c.ErrorLogger.Printf("failed to send peer request ID %d to %s: %v\n", reqID, n.Host, err)
		return nil, err
	}
	spr := new(ServerPeersSubscriptionResp)
	err = json.Unmarshal(resp, spr)
	if err != nil {
		c.ErrorLogger.Printf("error unmarshalling server peer subscription from %s req ID %d: %v\n", n.Host, reqID, err)
		return nil, err
	}
	peers, err := ParseServerPeersSubscriptionResp(spr)
	if err != nil {
		c.ErrorLogger.Printf("error parsing server peer subscription request response from %s for req ID %d: %v\n", n.Host, reqID, err)
		return nil, err
	}
	c.InfoLogger.Printf("successfully retrieved peer information from %s", n.Host)
	return peers, nil
}
