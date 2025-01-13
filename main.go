package main

import (
	"log"
	"net/http"
	"streaming-service/pkg/server"
)

func main() {
	streamServer := server.NewStreamServer()
	streamServer.Run()

	http.Handle("/stream", streamServer)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
