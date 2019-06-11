package metered_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/antonyho/go-project-demo/pkg/webserver/metered"
)

func ExampleNewWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	meteredWebServer := metered.NewWebServer(mux)
	meteredWebServer.Listen(":8080")
}

// TODO - Unit test functions on unexported functions
func TestNewWebServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	meteredWebServer := metered.NewWebServer(mux)
	go meteredWebServer.Listen(":9080")

	t.Run("GetRoot", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9080/")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if string(respBody) != "hello" {
			t.Errorf("Unexpected HTTP response content: %s. Expecting: hello.", string(respBody))
		}
	})
	t.Run("GetStat", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9080/stat")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if string(respBody) != "Total Request in Last 60 Seconds: 1" {
			t.Errorf("Unexpected HTTP response content: %s. Expecting: Total Request in Last 60 Seconds: 1.", string(respBody))
		}
	})
}
