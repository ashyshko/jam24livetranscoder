package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"jam24livetranscoder/protocol"

	"golang.org/x/net/websocket"
)

var wdpWrap WdpWrap

func Handler(ws *websocket.Conn) {
	session := newSession(ws, wdpWrap)
	defer session.Close()

	for {
		err := protocol.RecvServer(ws, &session)
		if err != nil {
			log.Printf("recv error: %s", err)
			return
		}
	}
}

func main() {
	wsPort := flag.Int("p", 8898, "TCP port for WS listen")
	flag.Parse()

	wdpWrap = newWdpWrap()

	http.Handle("/ws", websocket.Handler(Handler))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "ready")
	})

	http.HandleFunc("/preStop", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "not implemented...")
	})

	path := fmt.Sprintf(":%d", *wsPort)
	log.Printf("Listening on path %s...", path)
	err := http.ListenAndServe(path, nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %s", err)
	}
}
