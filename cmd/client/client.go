package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tylerchambers/electrumrelay/pkg/electrum"
	"io"
	"log"
	"net/http"
)

func main() {
	req := new(electrum.JSONRPCRequest)
	reqStr := fmt.Sprintf(`{"jsonrpc":"2.0","method":"blockchain.transaction.get","params":["d5845a4c59d7d3e86ab83650491ef2294552896599d036a440c08c52234e88f9", true],"id":0}`)
	err := json.Unmarshal([]byte(reqStr), req)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer([]byte(reqStr)))
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	fmt.Printf(string(b))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
