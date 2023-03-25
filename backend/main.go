package main

import (
	"backend/src/config"
	"backend/src/dispatcher_client"
	dispatcher_server "backend/src/dispatcher_client/dispathcer_server_proto"
	"backend/src/sse"
	server "backend/src/sse_server"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"message"`
}

func init() {
	config.BackendConfiguration = config.NewConfig(".env")
	sse.Clients = sse.NewSSE()
	server.GRPCServer = server.Start()
}

func main() {
	http.HandleFunc("/sse", func(writer http.ResponseWriter, request *http.Request) {
		sse.Clients.ServeHTTP(writer, request)
	})

	http.HandleFunc("/push-message", func(writer http.ResponseWriter, request *http.Request) {
		var message Message
		err := json.NewDecoder(request.Body).Decode(&message)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Message: %s", message.Text)

		conn, err := dispatcher_client.GetConnectionToDispatcherServer()

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer conn.Close()
		c := dispatcher_server.NewPushMessageClient(conn)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Printf(">>> Pushing message to dispatcher %s >>>", conn.Target())

		_, err = c.PushMessage(ctx, &dispatcher_server.PushMessageRequest{
			Message: message.Text,
		})

		if err != nil {
			log.Printf("Error pushing SSE: %s", err)
			return
		}
	})

	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		log.Fatal(err)
	}
}
