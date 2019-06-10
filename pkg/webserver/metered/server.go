// Package metered provides a metered web server
// which makes statistics on handled requests per minute.
package metered

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonyho/go-project-demo/pkg/webserver/framework"
)

const (
	// StatURL is used for handling HTTP request statistic inquiry.
	StatURL = "/stat"

	// HistoryFilename is the Request History Filename
	HistoryFilename = "req.his"
)

// webServer is a metered web server implementation.
type webServer struct {
	fw              framework.Wrapper // A web framework
	idleConnsClosed chan struct{}
	incReqChan      chan RequestInfo
	reqList         *list.List
	historyFile     *os.File
}

// NewWebServer returns a metered web server
// given the http.ServeMux to handle web requests.
func NewWebServer(mux *http.ServeMux) *webServer {
	webserver := new(webServer)
	mux.HandleFunc(StatURL, webserver.statistic) // This URL will be overrided
	f := framework.New(mux, webserver.filter, nil)
	webserver.fw = f
	return webserver
}

// filter gets the request information and passes it to a channel.
// The other side is a processing function.
func (s *webServer) filter(resp http.ResponseWriter, req *http.Request) {
	s.incReqChan <- RequestInfo{URL: req.RequestURI, Time: time.Now().UnixNano()}
}

// meter incoming web request
func (s *webServer) meter() {
	// Read history from file
	if err := s.restore(); err != nil {
		fmt.Println("Terminating web server due to history file retrieval failure.")
		panic(err)
	}

	for reqInfo := range s.incReqChan {
		// Store the request info to the list
		s.reqList.PushBack(reqInfo)
		// Persist to file
		if err := s.store(reqInfo); err != nil {
			// No specific requirement information about this minor error.
			// Tends to ignore this error in this example,
			// but the process should be terminated with error in real life.
		}
	}
}

// restore request information from file.
func (s *webServer) restore() error {
	windowStartPoint := time.Now().Add(-1 * time.Minute)
	buf := bufio.NewScanner(s.historyFile)
	for buf.Scan() {
		reqInfo := RequestInfo{}
		if err := json.Unmarshal(buf.Bytes(), &reqInfo); err != nil {
			log.Printf("Unmarshal request info from JSON error: %+v\n", err)
			return err
		}
		// Skip too old data for performance
		reqTime := time.Unix(0, reqInfo.Time)
		if reqTime.After(windowStartPoint) {
			// Store the history request info to the list
			s.reqList.PushBack(reqInfo)
		}
	}

	return nil
}

// store request information to file.
// This strategy is more robust then writing the list to a JSON file
// on every time to persist.
func (s *webServer) store(reqInfo RequestInfo) error {
	reqInfoJSON, err := json.Marshal(reqInfo)
	if err != nil {
		log.Printf("Marshal request info to JSON error: %+v\n", err)
		return err
	}
	_, err = fmt.Fprintln(s.historyFile, string(reqInfoJSON))
	if err != nil {
		log.Printf("Write to history file error: %+v\n", err)
		return err
	}
	return nil
}

// count the number of incoming request since last minute.
func (s *webServer) count() int {
	windowStartPoint := time.Now().Add(-1 * time.Minute)
	var next *list.Element
	total := 0
	for e := s.reqList.Front(); e != nil; e = next {
		reqTime := time.Unix(0, e.Value.(RequestInfo).Time)
		next = e.Next()
		if reqTime.Before(windowStartPoint) {
			// Discard expired request
			s.reqList.Remove(e)
		} else {
			// Count this request
			total++
		}
	}

	return total
}

// statistic handles HTTP request which inquiries statistics.
func (s *webServer) statistic(resp http.ResponseWriter, req *http.Request) {
	statFigure := s.count()
	fmt.Fprintf(resp, "Total Request in Last 60 Seconds: %d", statFigure)
}

// Listen listens listens on the TCP network address endpoint.
func (s *webServer) Listen(endpoint string) {
	s.incReqChan = make(chan RequestInfo)
	s.reqList = list.New()

	// Open history file from current working directory
	var err error
	s.historyFile, err = os.OpenFile(HistoryFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Open/Create history file error: %+v\n", err)
	}

	s.idleConnsClosed = make(chan struct{})
	httpSrv := &http.Server{
		Addr:    endpoint,
		Handler: s.fw,
	}
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan // Received interrupt signal

		// Shutdown the web server gracefully
		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %+v\n", err)
		}
		close(s.incReqChan)
		s.historyFile.Sync() // Write remaining request information to file
		s.historyFile.Close()

		close(s.idleConnsClosed)
	}()

	// Start the statistic meter
	go s.meter()

	// Start the HTTP Server
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %+v\n", err)
	}

	<-s.idleConnsClosed // Block until everything was gracefully terminated
}
