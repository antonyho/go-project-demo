// Package framework is a web server framework
// which has preprocess and postprocess filters for all requests.
package framework

import "net/http"

// Wrapper is an interface of a web server framework
// which provides preprocess and postprocess function filters.
type Wrapper interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	handle(http.ResponseWriter, *http.Request)
	preprocess(http.ResponseWriter, *http.Request)
	postprocess(http.ResponseWriter, *http.Request)
}

// Framework is a web server framework
// which provides preprocess and postprocess function filters.
type Framework struct {
	mux           *http.ServeMux
	preprocessor  func(http.ResponseWriter, *http.Request)
	postprocessor func(http.ResponseWriter, *http.Request)
}

// New returns a new Framework given optional handler, preprocessor,
// and postprocessor.
//
// If handler is not provided, a default http.ServeMux will be assigned.
func New(handler *http.ServeMux,
	preprocessor func(http.ResponseWriter, *http.Request),
	postprocessor func(http.ResponseWriter, *http.Request)) Wrapper {
	if handler == nil {
		handler = http.NewServeMux()
	}
	return &Framework{
		mux:           handler,
		preprocessor:  preprocessor,
		postprocessor: postprocessor,
	}
}

func (f Framework) handle(resp http.ResponseWriter, req *http.Request) {
	if f.mux == nil {
		f.mux = http.NewServeMux()
	}
	f.mux.ServeHTTP(resp, req)
}

func (f Framework) preprocess(resp http.ResponseWriter, req *http.Request) {
	if f.preprocessor != nil {
		f.preprocessor(resp, req)
	}
}

func (f Framework) postprocess(resp http.ResponseWriter, req *http.Request) {
	if f.postprocessor != nil {
		f.postprocessor(resp, req)
	}
}

// ServeHTTP dispatches the request like standard library web server
// with preprocessing and postprocessing.
func (f Framework) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	f.preprocess(resp, req)
	f.handle(resp, req)
	f.postprocess(resp, req)
}
