package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	opts, err := ParseOpts()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	addr := fmt.Sprintf(":%d", opts.port)
	server := NewServer(opts.commands)
	go server.Start()

	http.Handle("/ws", server.WebSocketHandler())
	if opts.webroot == "" {
		http.Handle("/", http.FileServer(assetFS()))
	} else {
		http.Handle("/", http.FileServer(http.Dir(opts.webroot)))
	}
	log.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
