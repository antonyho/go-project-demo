package metered

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewWebServer(t *testing.T) {
	// Remove the history file to make the test result idempotent
	if _, err := os.Stat(HistoryFilename); err == nil || !os.IsNotExist(err) {
		os.Remove(HistoryFilename)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "hello")
	})
	meteredWebServer := NewWebServer(mux)
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

func Test_webServer(t *testing.T) {
	server := NewWebServer(http.NewServeMux())
	server.incReqChan = make(chan RequestInfo, 1)
	defer close(server.incReqChan)
	testdataHistoryFilePath := fmt.Sprintf("testdata/%s", HistoryFilename)
	historyFile, err := os.OpenFile(testdataHistoryFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Skip(err)
	}
	defer historyFile.Close()
	server.historyFile = historyFile

	t.Run("restore", func(t *testing.T) {
		if err := server.restore(); err != nil {
			t.Error(err)
			t.Fail()
		}
	})
	t.Run("count", func(t *testing.T) {
		if server.count() != 3 {
			t.Fail()
		}
	})
	t.Run("statistic", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost/stat", nil)
		recorder := httptest.NewRecorder()
		server.statistic(recorder, req)
		resp := recorder.Result()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if string(respBody) != "Total Request in Last 60 Seconds: 3" {
			t.Logf("Unexcepted result: '%s'. Expected: 'Total Request in Last 60 Seconds: 3'\n", string(respBody))
			t.Fail()
		}
	})
	t.Run("filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost/stat", nil)
		recorder := httptest.NewRecorder()
		server.filter(recorder, req)

		if len(server.incReqChan) != 1 {
			t.Logf("Unexcepted channel length: %d. Expected: 1\n", len(server.incReqChan))
			t.Fail()
		}
	})
}
