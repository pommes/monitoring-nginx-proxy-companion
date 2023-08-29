package logparser

import "time"

type HttpRequest struct {
	Host           string
	SourceIp       string
	Timestamp      time.Time
	RequestMethod  string
	RequestPath    string
	HttpVersion    string
	HttpStatus     int
	BodyBytesSent  int
	HttpReferer    string
	UserAgent      string
	Browser        string
	BrowserVersion string
	Os             string
	Mobile         bool
	Country        string
	City           string
	Latency        float64
	xForwardedFor  string
	Bot            bool
}

type ILogParser interface {
	Parse(logLine string) (HttpRequest, error)
}
