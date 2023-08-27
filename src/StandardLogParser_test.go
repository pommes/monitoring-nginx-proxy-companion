package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParserLogLineWithoutLatency(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15"`

	userAgentParser := MssolaUserAgentParser{}
	mockIpLookupService := MockIpLookupService{}

	standardLogParser := StandardLogParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.sourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.requestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.requestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.httpVersion)
	assert.Equal(t, 200, httpRequest.httpStatus)
	assert.Equal(t, 20206, httpRequest.bodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.userAgent)
	assert.Equal(t, "Safari", httpRequest.browser)
	assert.Equal(t, "16.6", httpRequest.browserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.os)
	assert.Equal(t, false, httpRequest.mobile)
	assert.Equal(t, 0.0, httpRequest.latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}

func TestParserLogLineWithXForwardedFor(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 1.234`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "172.20.0.16:80"`

	userAgentParser := MssolaUserAgentParser{}
	mockIpLookupService := MockIpLookupService{}

	standardLogParser := StandardLogParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.sourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.requestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.requestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.httpVersion)
	assert.Equal(t, 200, httpRequest.httpStatus)
	assert.Equal(t, 20206, httpRequest.bodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.userAgent)
	assert.Equal(t, "Safari", httpRequest.browser)
	assert.Equal(t, "16.6", httpRequest.browserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.os)
	assert.Equal(t, false, httpRequest.mobile)
	assert.Equal(t, 0.0, httpRequest.latency)
	assert.Equal(t, "172.20.0.16:80", httpRequest.xForwardedFor)
}

func TestParserLogLineWithLatency(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 1.234`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" 1.234`

	userAgentParser := MssolaUserAgentParser{}
	mockIpLookupService := MockIpLookupService{}

	standardLogParser := StandardLogParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.sourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.requestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.requestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.httpVersion)
	assert.Equal(t, 200, httpRequest.httpStatus)
	assert.Equal(t, 20206, httpRequest.bodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.userAgent)
	assert.Equal(t, "Safari", httpRequest.browser)
	assert.Equal(t, "16.6", httpRequest.browserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.os)
	assert.Equal(t, false, httpRequest.mobile)
	assert.Equal(t, 1.234, httpRequest.latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}
