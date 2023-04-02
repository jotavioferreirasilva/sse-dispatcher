package test

import (
	"backend/src/handlers"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestWhenPushingAMessageThenAllClientsShouldReceiveIt(t *testing.T) {
	client := &http.Client{}
	var wg sync.WaitGroup

	reqClient1, _ := http.NewRequest(http.MethodGet, "http://localhost/sse?client=1", nil)
	reqClient1.Header.Set("Content-Type", "text/event-stream")
	reqClient1.Header.Set("Connection", "keep-alive")

	wg.Add(1)
	go func() {
		responseClient1, _ := client.Do(reqClient1)
		assert.Equal(t, http.StatusOK, responseClient1.StatusCode)
		body := make([]byte, 18)
		io.ReadFull(responseClient1.Body, body)
		assert.Equal(t, "data: Hello world!", string(body))
		wg.Done()
	}()

	reqClient2, _ := http.NewRequest(http.MethodGet, "http://localhost/sse?client=2", nil)
	reqClient2.Header.Set("Content-Type", "text/event-stream")
	reqClient2.Header.Set("Connection", "keep-alive")

	wg.Add(1)
	go func() {
		responseClient2, _ := client.Do(reqClient2)
		assert.Equal(t, http.StatusOK, responseClient2.StatusCode)
		body := make([]byte, 18)
		io.ReadFull(responseClient2.Body, body)
		assert.Equal(t, "data: Hello world!", string(body))
		wg.Done()
	}()

	message := handlers.Message{Text: "Hello world!"}
	body, _ := json.Marshal(message)

	requestMessage, err := http.NewRequest("POST", "http://localhost/push-message", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	client.Do(requestMessage)

	wg.Wait()
}
