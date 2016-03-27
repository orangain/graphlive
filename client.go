package main

import (
	"golang.org/x/net/websocket"
)

type Client struct {
	id int
	ws *websocket.Conn
	// channels
	removeClientCh chan *Client
}

func NewClient(ws *websocket.Conn, removeClientCh chan *Client) *Client {
	return &Client{0, ws, removeClientCh}
}

func (c *Client) Start() {
	for {
		var message string
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			c.removeClientCh <- c
			return
		}
	}
}

func (c *Client) Send(metrics *Metrics) {
	websocket.JSON.Send(c.ws, metrics)
}
