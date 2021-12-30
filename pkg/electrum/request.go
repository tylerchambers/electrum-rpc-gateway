package electrum

import (
	"encoding/json"
	"fmt"
	"net"
)

type JSONRPCRequest struct {
	Version string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (r *JSONRPCRequest) Send(conn net.Conn) error {
	reqBytes, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(conn, "%s\n", reqBytes)
	if err != nil {
		return err
	}
	return nil
}
