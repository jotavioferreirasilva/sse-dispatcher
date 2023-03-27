package main

import (
	"backend/src/config"
	"backend/src/handlers"
	"backend/src/sse"
	server "backend/src/sse_server"
	"log"
	"net/http"
)

func init() {
	config.BackendConfiguration = config.NewConfig(".env")
	sse.Clients = sse.NewSSE()
	server.GRPCServer = server.Start()
}

func main() {
	http.HandleFunc("/sse", handlers.SSE)
	http.HandleFunc("/push-message", handlers.PushMessage)

	log.Printf("Server started on port %d", 5000)
	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		log.Fatal(err)
	}
}
