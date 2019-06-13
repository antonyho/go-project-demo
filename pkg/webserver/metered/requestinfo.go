package metered

import (
	"container/list"
	"log"
	"sync"
	"time"
)

// RequestInfo stores the basic information of an HTTP request.
// URL is the request URL.
// Request time is the server's local time down to nano-second precision in Unix epoch.
type RequestInfo struct {
	URL  string `json:"url"`
	Time int64  `json:"request_time"`
}

// RequestInfoPool is a sychronised pool which stores RequestInfo.
// It is thread safe.
type RequestInfoPool struct {
	pool *list.List
	m    sync.Mutex
}

// NewRequestInfoPool returns a RequestInfoPool.
func NewRequestInfoPool() *RequestInfoPool {
	return &RequestInfoPool{
		pool: new(list.List),
		m:    sync.Mutex{},
	}
}

// Add any number of RequestInfo into the pool.
func (p *RequestInfoPool) Add(requestInfos ...RequestInfo) {
	p.m.Lock()
	defer p.m.Unlock()
	for _, reqInfo := range requestInfos {
		log.Printf("%+v\n", reqInfo)
		p.pool.PushBack(reqInfo)
	}
}

// ClearRecord any RequestInfo from the pool which is earlier than before.
func (p *RequestInfoPool) ClearRecord(before time.Time) {
	p.m.Lock()
	defer p.m.Unlock()
	var next *list.Element
	for e := p.pool.Front(); e != nil; e = next {
		reqTime := time.Unix(0, e.Value.(RequestInfo).Time)
		next = e.Next()
		if reqTime.Before(before) {
			p.pool.Remove(e)
		}
	}
}

// Count the number of RequestInfo in the pool.
func (p *RequestInfoPool) Count() int {
	p.m.Lock()
	defer p.m.Unlock()
	return p.pool.Len()
}
