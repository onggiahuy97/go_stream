package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type StreamServer struct {
	clients    map[chan []byte]bool
	broadcast  chan []byte
	register   chan chan []byte
	unregister chan chan []byte
}

func NewStreamServer() *StreamServer {
	return &StreamServer{
		clients:    make(map[chan []byte]bool),
		broadcast:  make(chan []byte),
		register:   make(chan chan []byte),
		unregister: make(chan chan []byte),
	}
}

func (s *StreamServer) Run() {
	go func() {
		for {
			select {
			case client := <-s.register:
				s.clients[client] = true
				log.Printf("Client registered. Total clients: %d", len(s.clients))
			case client := <-s.unregister:
				if _, ok := s.clients[client]; ok {
					delete(s.clients, client)
					close(client)
					log.Printf("Client unregistered. Total clients: %d", len(s.clients))
				}
			case message := <-s.broadcast:
				for client := range s.clients {
					select {
					case client <- message:
					default:
						close(client)
						delete(s.clients, client)
					}
				}
			}
		}
	}()

	// Simulate sending data every second
	go func() {
		for {
			message := []byte(fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC3339)))
			s.broadcast <- message
			time.Sleep(time.Second)
		}
	}()
}

func (s *StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)
	s.register <- messageChan

	notify := r.Context().Done()
	go func() {
		<-notify
		s.unregister <- messageChan
	}()

	for {
		select {
		case <-notify:
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
