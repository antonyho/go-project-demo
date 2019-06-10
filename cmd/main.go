package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/antonyho/go-project-demo/pkg/webserver/metered"
)

var (
	hostname string
	port     int
)

func init() {
	flag.StringVar(&hostname, "h", "", "Hostname for the metered web server")
	flag.IntVar(&port, "p", 80, "Port to listen for the metered web server")
	flag.Parse()
}

func main() {
	endpoint := fmt.Sprintf("%s:%d", hostname, port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	meteredWebServer := metered.NewWebServer(mux)
	meteredWebServer.Listen(endpoint)
}
