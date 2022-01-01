package main

import (
	"fmt"
	"github.com/tylerchambers/electrumrelay/pkg/electrum"
	"github.com/tylerchambers/electrumrelay/pkg/relay"
	"log"
)

func main() {
	initialNode := &electrum.Node{
		Host:         "electrum.blockstream.info",
		IP:           "",
		Version:      "",
		SSLPort:      50002,
		TCPPort:      0,
		PruningLimit: 0,
	}
	r := new(relay.Relay)
	err := r.Init(initialNode)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range r.Peers {
		fmt.Println(v)
	}

}
