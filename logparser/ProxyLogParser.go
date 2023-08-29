package logparser

import (
	"errors"
	"fmt"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/useragentparser"
	"regexp"
	"strconv"
)

type ProxyLogParser struct {
	UserAgentParser useragentparser.IUserAgentParser
	IpLookupService geoip.IIpLookupService
}

func (logParser ProxyLogParser) Parse(logLine string) (HttpRequest, error) {
	const LogLineProxyRegex = `^(?P<hostname>\S+) (?P<remote_addr>\S+) - (?P<remote_user>\S+) \[(?P<time_local>.+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<body_bytes_sent>\S+) "(?P<referer>[^\"]*)" "(?P<user_agent>[^\"]*)" "(?P<http_x_forwarded_for>[^\"]*)"`
	var logPattern = regexp.MustCompile(LogLineProxyRegex)
	match := logPattern.FindStringSubmatch(logLine)
	if len(match) <= 0 {
		return HttpRequest{}, errors.New(fmt.Sprintf("Log line did not match nginx-proxy log format: '%s'", logLine))
	}

	result := make(map[string]string)
	for i, name := range logPattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	//fmt.Printf("%+v\n", result)

	httpRequest := HttpRequest{}
	httpRequest.Host = result["hostname"]
	httpRequest.SourceIp = result["remote_addr"]
	httpRequest.Timestamp = convertDateStringToTime(result["time_local"])
	httpRequest.RequestMethod = result["method"]
	httpRequest.RequestPath = result["path"]
	httpRequest.HttpVersion = result["protocol"]
	httpRequest.HttpStatus, _ = strconv.Atoi(result["status"])
	httpRequest.BodyBytesSent, _ = strconv.Atoi(result["body_bytes_sent"])
	httpRequest.HttpReferer = result["referer"]
	httpRequest.UserAgent = result["user_agent"]
	httpRequest.Latency = 0.0
	httpRequest.xForwardedFor = result["http_x_forwarded_for"]

	parseUserAgentAndSetFields(logParser.UserAgentParser, httpRequest.UserAgent, &httpRequest)
	lookupIpAndSetFields(logParser.IpLookupService, httpRequest.SourceIp, &httpRequest)

	return httpRequest, nil
}
