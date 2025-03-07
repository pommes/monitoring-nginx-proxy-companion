package logline

import (
	"github.com/stretchr/testify/assert"
	"log"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/useragent"
	"testing"
	"time"
)

func TestStandardLogParserLogLineWithoutLatency(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15"`

	userAgentParser := useragent.MssolaParser{}
	mockIpLookupService := geoip.IPLocatorMock{}

	standardLogParser := StandardParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.Host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.SourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.Timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.RequestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.RequestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.HttpVersion)
	assert.Equal(t, 200, httpRequest.HttpStatus)
	assert.Equal(t, 20206, httpRequest.BodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.UserAgent)
	assert.Equal(t, "Safari", httpRequest.Browser)
	assert.Equal(t, "16.6", httpRequest.BrowserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.Os)
	assert.Equal(t, false, httpRequest.Mobile)
	assert.Equal(t, 0.0, httpRequest.Latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}

func TestStandardLogParserLogLineErrorInProduction1(t *testing.T) {
	t.Skip("This parser does not support this case.")
	line := `3.89.123.261 22.102.114.111 - - [29/Aug/2023:07:47:30 +0000] "GET /favicon.ico HTTP/1.1" 502 552 "" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11" "172.20.0.14:8007"`
	//line := `fh-warzone.de 63.143.42.253 - - [27/Aug/2023:20:37:53 +0000] "HEAD /forum/      HTTP/1.1" 200 0   "http://fh-warzone.de" "Mozilla/5.0+(compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)" "172.20.0.25:80"`

	userAgentParser := useragent.MssolaParser{}
	mockIpLookupService := geoip.IPLocatorMock{}

	standardLogParser := StandardParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		log.Println("ERROR: Line does not match:", line)
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "3.89.123.261", httpRequest.Host)
	assert.Equal(t, "22.102.114.111", httpRequest.SourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.Timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.RequestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.RequestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.HttpVersion)
	assert.Equal(t, 200, httpRequest.HttpStatus)
	assert.Equal(t, 20206, httpRequest.BodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.UserAgent)
	assert.Equal(t, "Safari", httpRequest.Browser)
	assert.Equal(t, "16.6", httpRequest.BrowserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.Os)
	assert.Equal(t, false, httpRequest.Mobile)
	assert.Equal(t, 0.0, httpRequest.Latency)
	assert.Equal(t, "172.20.0.16:80", httpRequest.xForwardedFor)
}

func TestStandardLogParserLogLineWithXForwardedFor(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 1.234`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "172.20.0.16:80"`
	//line := `fh-warzone.de 63.143.42.253 - - [27/Aug/2023:20:37:53 +0000] "HEAD /forum/ HTTP/1.1" 200 0 "http://fh-warzone.de" "Mozilla/5.0+(compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)" "172.20.0.25:80"`
	//line := `devarcs.com 54.201.119.245 - - [27/Aug/2023:20:36:30 +0000] "GET / HTTP/1.1" 301 162 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36" "-"`
	//line := `zahlensender.net 162.55.94.150 - - [27/Aug/2023:21:11:29 +0000] "GET /impuls/feed/m4a/ HTTP/2.0" 304 0 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36" "172.20.0.16:80"`

	userAgentParser := useragent.MssolaParser{}
	mockIpLookupService := geoip.IPLocatorMock{}

	standardLogParser := StandardParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.Host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.SourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.Timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.RequestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.RequestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.HttpVersion)
	assert.Equal(t, 200, httpRequest.HttpStatus)
	assert.Equal(t, 20206, httpRequest.BodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.UserAgent)
	assert.Equal(t, "Safari", httpRequest.Browser)
	assert.Equal(t, "16.6", httpRequest.BrowserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.Os)
	assert.Equal(t, false, httpRequest.Mobile)
	assert.Equal(t, 0.0, httpRequest.Latency)
	assert.Equal(t, "172.20.0.16:80", httpRequest.xForwardedFor)
}

func TestStandardLogParserLogLineWithLatency(t *testing.T) {
	//line := `blog.kroepfl.io 193.80.91.32 - - [27/May/2017:19:26:27 +0000] "GET /wp-content/uploads/2017/04/Untitled.png HTTP/1.1" 404 18000 "https://blog.kroepfl.io/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 1.234`
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "1.234""`

	userAgentParser := useragent.MssolaParser{}
	mockIpLookupService := geoip.IPLocatorMock{}

	standardLogParser := StandardParser{
		userAgentParser: userAgentParser,
		ipLookupService: mockIpLookupService,
	}

	httpRequest, err := standardLogParser.Parse(line)

	if err != nil {
		t.Fail()
	}

	ti := time.Unix(1495913187, 0)

	assert.Equal(t, "zahlensender.net", httpRequest.Host)
	assert.Equal(t, "2aff:e202:3001:1854:e::1", httpRequest.SourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.Timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.RequestMethod)
	assert.Equal(t, "/path/data/my-data.html", httpRequest.RequestPath)
	assert.Equal(t, "HTTP/2.0", httpRequest.HttpVersion)
	assert.Equal(t, 200, httpRequest.HttpStatus)
	assert.Equal(t, 20206, httpRequest.BodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15", httpRequest.UserAgent)
	assert.Equal(t, "Safari", httpRequest.Browser)
	assert.Equal(t, "16.6", httpRequest.BrowserVersion)
	assert.Equal(t, "Intel Mac OS X 10_15_7", httpRequest.Os)
	assert.Equal(t, false, httpRequest.Mobile)
	assert.Equal(t, 1.234, httpRequest.Latency)
	assert.Equal(t, "", httpRequest.xForwardedFor)
}
