package sse

import (
	"fmt"
	"net/http"
	"sync"
)

var Clients Client

type Client struct {
	Client map[string]chan string
	sync.Mutex
}

func NewSSE() Client {
	var sseClient Client
	sseClient.Client = make(map[string]chan string)
	return sseClient
}

func (sse *Client) AddClient(client string, messageChannel chan string) {
	sse.Lock()
	defer sse.Unlock()
	sse.Client[client] = messageChannel
}

func (sse *Client) RemoveClient(client string) {
	sse.Lock()
	defer sse.Unlock()
	delete(sse.Client, client)
}

func (sse *Client) GetClient(client string) (chan string, bool) {
	sse.Lock()
	defer sse.Unlock()
	messageChannel, ok := sse.Client[client]
	return messageChannel, ok
}

func (sse *Client) Publish(message string) {
	sse.Lock()
	defer sse.Unlock()
	for _, messageChannel := range sse.Client {
		messageChannel <- message
	}
}

func (sse *Client) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	clientID := request.URL.Query().Get("client")
	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	messageChannel := make(chan string)
	sse.AddClient(clientID, messageChannel)
	var wg sync.WaitGroup

	defer func() {
		sse.RemoveClient(clientID)
		close(messageChannel)
	}()

	wg.Add(1)
	go func() {
		for {
			select {
			case message, ok := <-messageChannel:
				if !ok {
					return
				}
				fmt.Fprintf(writer, "data: %s\n\n", message)
				flusher.Flush()
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
