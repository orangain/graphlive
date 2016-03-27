package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":9999"
	commands := []string{
		"uptime | awk '{print $10}'",
		"python -c 'import random; print(random.random())'",
	}

	server := NewServer(commands)
	go server.Start()

	http.Handle("/ws", server.WebSocketHandler())
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
