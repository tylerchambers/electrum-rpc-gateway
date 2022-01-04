package main

import (
	"fmt"
	"github.com/tylerchambers/electrumrelay/pkg/electrum"
	"github.com/tylerchambers/electrumrelay/pkg/relay"
	"log"
	"net/http"
)

type server struct {
	relay  *relay.Relay
	router *http.ServeMux
}

func main() {
	s := server{
		router: http.NewServeMux(),
	}

	s.router.HandleFunc("/", s.handleRelay)
	initialNode := &electrum.Node{
		Host:         "electrum.blockstream.info",
		IP:           "",
		Version:      "",
		SSLPort:      50002,
		TCPPort:      0,
		PruningLimit: 0,
	}

	// set up the relay and register initial peers
	ec := electrum.NewClient(log.Default(), log.Default(), log.Default())
	r := relay.NewRelay([]electrum.Node{}, []string{}, ec)
	err := r.Bootstrap(initialNode)
	if err != nil {
		log.Fatal(err)
	}

	s.relay = r

	for _, v := range r.Peers {
		fmt.Println(v)
	}

	log.Fatal(http.ListenAndServe(":8080", s.router))

}

func (s *server) handleRelay(w http.ResponseWriter, r *http.Request) {
	req, err := s.relay.ValidateRequest(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte("a error, see logs for details"))
		return
	}
	resp, err := s.relay.ForwardRequest(req)
	if err != nil {
		w.Write([]byte("b error, see logs for details"))
		log.Println(err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		w.Write([]byte("c error, see logs for details"))
		log.Println(err)
		return
	}
}
