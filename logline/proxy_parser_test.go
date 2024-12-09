package logline

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/useragent"
	"regexp"
	"testing"
	"time"
)

var (
	userAgentParser     = useragent.MssolaParser{}
	mockIpLookupService = geoip.IPLocatorMock{}

	logParser = ProxyParser{
		UserAgentParser: userAgentParser,
		IPLocator:       mockIpLookupService,
	}
)

func failNowIfErr(t *testing.T, err error, line string) {
	if err != nil {
		t.Log(fmt.Sprintf("ERROR: Unexpected error: '%s'\n", line))
		t.FailNow()
	}
}

func failNowIfNoErr(t *testing.T, err error, line string) {
	if err == nil {
		t.Log(fmt.Sprintf("ERROR: Did expect error, but passed: '%s'\n", line))
		t.FailNow()
	}
}

func TestProxyLogParserLogLineErrorInProduction1(t *testing.T) {

	line := `3.89.123.261 22.102.114.111 - - [29/Aug/2023:07:47:30 +0000] "GET /favicon.ico HTTP/1.1" 502 552 "" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11" "172.20.0.14:8007"`
	//line := `fh-warzone.de 63.143.42.253 - - [27/Aug/2023:20:37:53 +0000] "HEAD /forum/      HTTP/1.1" 200 0   "http://fh-warzone.de" "Mozilla/5.0+(compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)" "172.20.0.25:80"`

	httpRequest, err := logParser.Parse(line)
	failNowIfErr(t, err, line)

	ti := time.Unix(1693295250, 0)

	assert.Equal(t, "3.89.123.261", httpRequest.Host)
	assert.Equal(t, "22.102.114.111", httpRequest.SourceIp)
	assert.Equal(t, ti.Unix(), httpRequest.Timestamp.Unix())
	assert.Equal(t, "GET", httpRequest.RequestMethod)
	assert.Equal(t, "/favicon.ico", httpRequest.RequestPath)
	assert.Equal(t, "HTTP/1.1", httpRequest.HttpVersion)
	assert.Equal(t, 502, httpRequest.HttpStatus)
	assert.Equal(t, 552, httpRequest.BodyBytesSent)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11", httpRequest.UserAgent)
	assert.Equal(t, "Chrome", httpRequest.Browser)
	assert.Equal(t, "17.0.963.56", httpRequest.BrowserVersion)
	assert.Equal(t, "Intel Mac OS X 10_7_0", httpRequest.Os)
	assert.Equal(t, false, httpRequest.Mobile)
	assert.Equal(t, 0.0, httpRequest.Latency)
	assert.Equal(t, "172.20.0.14:8007", httpRequest.xForwardedFor)
}

func TestProxyLogParserLogLineErrorInProduction2(t *testing.T) {

	line := `awstats.tyranus.de 181.214.164.109 - - [30/Aug/2023:08:07:07 +0000] ":\x92T\x05ib\xE8\x0Ek_V\x08\xDD=x\xAB\xC2\x13\x22\xB88\x1B\x01\x07\xA6\xB1~\xE0Ap\x8D\x96\xF3 \xB9\xDB\x0CEN#5h[\xE4\xC5\x16\xF7wBr=\xB1" 400 150 "-" "-" "-"`

	_, err := logParser.Parse(line)
	failNowIfNoErr(t, err, line) // only 1 whitespace found in block "<request-method> <path> <http-version>"
}

func TestProxyLogParserLogLineWithXForwardedFor(t *testing.T) {
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "172.20.0.16:80"`

	httpRequest, err := logParser.Parse(line)
	failNowIfErr(t, err, line)

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

func TestProxyLogParserLogLineNoRegexMatch(t *testing.T) {
	line := `zahlensender.net 2aff:e202:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" 1.234`

	_, err := logParser.Parse(line)
	failNowIfNoErr(t, err, line)
}

func TestProxyLogLogParserImplementation(t *testing.T) {
	const logEntry = `abc.net 0eaa:e103:3001:1854:e::1 - - [27/May/2017:19:26:27 +0000] "GET /path/data/my-data.html HTTP/2.0" 200 20206 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15" "172.20.0.16:80"`
	var combinedLogPattern = regexp.MustCompile(`^(?P<hostname>\S+) (?P<remote_addr>\S+) - (?P<remote_user>\S+) \[(?P<time_local>.+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<body_bytes_sent>\S+) "(?P<referer>[^\"]*)" "(?P<user_agent>[^\"]*)" "(?P<http_x_forwarded_for>[^\"]*)"`)

	match := combinedLogPattern.FindStringSubmatch(logEntry)
	result := make(map[string]string)
	for i, name := range combinedLogPattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	fmt.Printf("%+v\n", result)
}

func TestProxyLogRemovesANSICodesFromHostname(t *testing.T) {
	line := `[0mprom-grafana.tyranus.de 127.0.0.1 - - [29/Aug/2023:07:47:30 +0000] "GET / HTTP/1.1" 200 12345 "-" "Mozilla/5.0" "-"`

	httpRequest, err := logParser.Parse(line)
	failNowIfErr(t, err, line)

	expectedHost := "prom-grafana.tyranus.de"
	assert.Equal(t, expectedHost, httpRequest.Host, "The hostname should have ANSI codes removed")
}
