package goat

import (
	"net"
	"net/http"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(l net.Listener) bool
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle incoming HTTP connections and serve
func (h HttpConnHandler) Handle(l net.Listener, logChan chan string) bool {
	http.HandleFunc("/announce", parseHttp)

	err := http.Serve(l, nil)
	if err != nil {
		logChan <- err.Error()
	}

	return true
}

// Parse incoming HTTP connections before making tracker calls
func parseHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", APP+"-git")
	w.Write([]byte("announce successful"))
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l net.Listener) bool {
	return true
}
