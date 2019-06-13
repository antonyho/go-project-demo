package metered

import (
	"testing"
	"time"
)

func TestRequestInfoPool(t *testing.T) {
	var pool *RequestInfoPool
	var pause time.Time
	t.Run("NewRequestInfoPool", func(t *testing.T) {
		pool = NewRequestInfoPool()
		if pool.Count() != 0 {
			t.Error("NewRequestInfoPool initialisation failed")
		}
	})
	t.Run("Add", func(t *testing.T) {
		reqInfos := make([]RequestInfo, 10)
		for i := 0; i < 10; i++ {
			reqInfos = append(reqInfos, RequestInfo{URL: "/", Time: time.Now().UnixNano()})
		}
		pause = time.Now()
		pool.Add(reqInfos...)

		reqInfos = make([]RequestInfo, 10)
		for i := 0; i < 10; i++ {
			reqInfos = append(reqInfos, RequestInfo{URL: "/", Time: time.Now().UnixNano()})
		}
		pool.Add(reqInfos...)
	})
	t.Run("ClearRecord", func(t *testing.T) {
		pool.ClearRecord(pause)
	})
	t.Run("Count", func(t *testing.T) {
		count := pool.Count()
		if count != 10 {
			t.Errorf("Unexpected count value: %d. Expected: 10.", count)
		}
	})
}
