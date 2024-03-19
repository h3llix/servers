package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func connectWebSocket(url string, message string) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Send message to WebSocket server
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		fmt.Println("Error sending message to WebSocket:", err)
		return
	}

	// Read response from WebSocket server
	_, response, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error reading response from WebSocket:", err)
		return
	}

	fmt.Println("Response from", url, ":", string(response), "\n")
}

func sendHTTPRequest(url, message string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "text/plain", strings.NewReader(message))
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response from", url, ":", string(body), "\n")

}

func main() {
	http1URL := "http://localhost:8080/http1"
	sendHTTPRequest(http1URL, "Hello from HTTP/1.1 client")

	websocketURL := "ws://localhost:8080/websocket"
	connectWebSocket(websocketURL, "Hello from WebSocket client")

	http2URL := "https://localhost:8081/http2"
	sendHTTPRequest(http2URL, "Hello from HTTP/2 client")
}
