package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	httpAddr  = flag.String("addr", ":8080", "ip:port")
	httpsAddr = flag.String("tls addr", ":8081", "ip:port")
)

func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func http1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URI path: %s\n", r.URL.Path)

	fmt.Fprintf(w, "Headers:\n")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		fmt.Fprintf(w, "\nBody:\n")
		buf := make([]byte, 1024)
		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				fmt.Fprintf(w, "%s", buf[:n])
			}
			if err != nil {
				break
			}
		}
	}
}

func http2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URI path: %s\n", r.URL.Path)

	fmt.Fprintf(w, "Headers:\n")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		fmt.Fprintf(w, "\nBody:\n")
		buf := make([]byte, 1024)
		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				fmt.Fprintf(w, "%s", buf[:n])
			}
			if err != nil {
				break
			}
		}
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection")
		return
	}
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current directory")
	}
	commnand := fmt.Sprintf("%s/%s", cwd, "script.sh")
	cmd := exec.Command(commnand)
	cmd.Run()

	http_mux := mux.NewRouter()
	http_mux.HandleFunc("/http1", http1Handler)
	http_mux.HandleFunc("/websocket", websocketHandler)
	go http.ListenAndServe(*httpAddr, http_mux)
	fmt.Printf("starting http and ws server on %s\n", *httpAddr)

	https_mux := mux.NewRouter()
	https_mux.HandleFunc("/http2", http2Handler)
	go http.ListenAndServeTLS(*httpsAddr, "domain.crt", "domain.key", https_mux)

	fmt.Printf("starting https server on %s\n", *httpsAddr)

	handleSignals()
}
