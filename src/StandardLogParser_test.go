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
	assert.Equal(t, 0, httpRequest.latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}

func TestParserLogLineWithXForwardedFor(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 1.234`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "172.20.0.16:80"`
	//line := `fh-warzone.de 63.143.42.253 - - [27/Aug/2023:20:37:53 +0000] "HEAD /forum/ HTTP/1.1" 200 0 "http://fh-warzone.de" "Mozilla/5.0+(compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)" "172.20.0.25:80"`
	//line := `devarcs.com 54.201.119.245 - - [27/Aug/2023:20:36:30 +0000] "GET / HTTP/1.1" 301 162 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36" "-"`
	//line := `zahlensender.net 162.55.94.150 - - [27/Aug/2023:21:11:29 +0000] "GET /impuls/feed/m4a/ HTTP/2.0" 304 0 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36" "172.20.0.16:80"`

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
	assert.Equal(t, 0, httpRequest.latency)
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
	assert.Equal(t, 1, httpRequest.latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}
