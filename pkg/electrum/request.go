package electrum

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// JSONRPCRequest represents a JSON RPC Request.
type JSONRPCRequest struct {
	Version string        `json:"jsonrpc"`
	ID      int           `json:"id,omitempty"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// NewJSONRPCRequest creates a new JSONRPCRequest
func NewJSONRPCRequest(version string, ID int, method string, params []interface{}) *JSONRPCRequest {
	return &JSONRPCRequest{Version: version, ID: ID, Method: method, Params: params}
}

// Send sends the JSONRPCRequest to the specified conn.
func (r *JSONRPCRequest) Send(conn net.Conn) ([]byte, error) {
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return nil, err
	}
	reqBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(conn, "%s\n", reqBytes)
	if err != nil {
		return nil, err
	}
	resp, err := bufio.NewReader(conn).ReadBytes(byte('\n'))
	if err != nil {
		return nil, err
	}
	return resp, err
}
