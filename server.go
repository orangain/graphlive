package main

import (
	"golang.org/x/net/websocket"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	Value float64 `json:"value"`
	Label string  `json:"label"`
	index int
}

type Metrics struct {
	Time    int      `json:"time"`
	Metrics []Metric `json:"metrics"`
}

type Server struct {
	nextClientId int
	clients      map[int]*Client
	commands     []string
	// channels
	addClientCh    chan *Client
	removeClientCh chan *Client
	sendMetricsCh  chan *Metrics
}

func NewServer(commands []string) *Server {
	return &Server{
		0,
		make(map[int]*Client),
		commands,
		make(chan *Client),
		make(chan *Client),
		make(chan *Metrics),
	}
}

func (s *Server) Start() {
	// Collect metrics every seconds
	go func() {
		for {
			s.CollectMetrics()
			<-time.After(time.Second)
		}
	}()

	// Synchrnize channels
	for {
		select {
		case c := <-s.addClientCh:
			s.AddClient(c)
		case c := <-s.removeClientCh:
			s.RemoveClient(c)
		case m := <-s.sendMetricsCh:
			s.SendMetrics(m)
		}
	}
}

func (s *Server) AddClient(c *Client) {
	c.id = s.nextClientId
	s.clients[c.id] = c
	s.nextClientId++
	log.Println("Added a new client.")
	log.Println("Now", len(s.clients), "clients are connected.")
}

func (s *Server) RemoveClient(c *Client) {
	delete(s.clients, c.id)
	log.Println("Removed a client.")
	log.Println("Now", len(s.clients), "clients are connected.")
}

func (s *Server) CollectMetrics() {
	unixMilli := int(time.Now().UnixNano() / 1000000)
	metrics := Metrics{unixMilli, make([]Metric, len(s.commands))}

	metricCh := make(chan *Metric)

	for i, c := range s.commands {
		go func(index int, command string) {
			out, _ := exec.Command("sh", "-c", command).Output()
			val, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
			metricCh <- &Metric{val, command, index}
		}(i, c)
	}

	// Wait all metrics
	for i := 0; i < len(s.commands); i++ {
		m := <-metricCh
		metrics.Metrics[m.index] = *m
	}

	for i, m := range metrics.Metrics {
		log.Printf("[%v] %v - %v\n", i, m.Value, m.Label)
	}
	s.sendMetricsCh <- &metrics
}

func (s *Server) SendMetrics(metrics *Metrics) {
	for _, c := range s.clients {
		c.Send(metrics)
	}
}

func (s *Server) WebSocketHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		c := NewClient(ws, s.removeClientCh)
		s.addClientCh <- c
		c.Start()
	})
}
