package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tylerchambers/electrumrelay/pkg/electrum"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "de.poiuty.com:50002")
	if err != nil {
		fmt.Println("1")
		panic(err)
	}
	log.Println("connected!")
	defer conn.Close()
	m := new(electrum.JSONRPCRequest)
	req := `{"jsonrpc":"2.0","method":"server.peers.subscribe","params":[],"id":64}`
	json.Unmarshal([]byte(req), m)
	fmt.Println(m)
	if err != nil {
		fmt.Println("2")
		panic(err)
	}
	reqbytes, err := json.Marshal(m)
	_, err = fmt.Fprintf(conn, "%s\n", reqbytes)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		fmt.Println("3")
		log.Fatal(err)
	}
	log.Println("request sent!")
	respBytes, _ := bufio.NewReader(conn).ReadBytes(byte('\n'))
	fmt.Println("got bytes:")
	fmt.Println(string(respBytes))
	resp := new(electrum.ServerPeersSubscriptionResp)
	fmt.Println("unmarshalling")
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response successfully marshalled...")
	fmt.Print(resp)
	fmt.Println("Parsing the response...")
	parsed, err := electrum.ParseServerPeersSubscriptionResp(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(parsed)
}
