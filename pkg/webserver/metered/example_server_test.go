package metered_test

import (
	"fmt"
	"net/http"

	"github.com/antonyho/go-project-demo/pkg/webserver/metered"
)

func ExampleNewWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	meteredWebServer := metered.NewWebServer(mux)
	go meteredWebServer.Listen(":8080")
}
