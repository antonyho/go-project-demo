package framework_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/antonyho/go-project-demo/pkg/webserver/framework"
)

// TestFramework tests the web server framework with filters.
// It assumes your ports 9081, 9082, 9083 are all free.
func TestFramework(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	mux.HandleFunc("/world", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "world")
	})
	pre := func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "pre")
	}
	post := func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "post")
	}
	client := new(http.Client)

	t.Run("NoHandlerAndFilter", func(t *testing.T) {
		f := framework.Framework{}
		go http.ListenAndServe(":9081", f)

		resp, err := client.Get("http://localhost:9081/hello")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if resp.StatusCode != 404 {
			t.Errorf("Unexpected HTTP response: %d. Expecting: 404.", resp.StatusCode)
		}
	})
	t.Run("NoFilter", func(t *testing.T) {
		f := framework.New(mux, nil, nil)
		go http.ListenAndServe(":9082", f)

		resp, err := client.Get("http://localhost:9082/hello")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		defer resp.Body.Close()
		if string(respBody) != "hello" {
			t.Errorf("Unexpected response body: %s. Expecting: hello", string(respBody))
		}
	})
	t.Run("WithFilter", func(t *testing.T) {
		f := framework.New(mux, pre, post)
		go http.ListenAndServe(":9083", f)

		resp, err := client.Get("http://localhost:9083/hello")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		defer resp.Body.Close()
		if string(respBody) != "prehellopost" {
			t.Errorf("Unexpected response body: %s. Expecting: prehellopost", string(respBody))
		}
	})
}
