package metered

// RequestInfo stores the basic information of an HTTP request.
// URL is the request URL.
// Request time is the server's local time down to nano-second precision in Unix epoch.
type RequestInfo struct {
	URL  string `json:"url"`
	Time int64  `json:"request_time"`
}
