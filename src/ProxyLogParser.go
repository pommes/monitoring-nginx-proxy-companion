package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ProxyLogParser struct {
	userAgentParser IUserAgentParser
	ipLookupService IIpLookupService
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
	httpRequest.host = result["hostname"]
	httpRequest.sourceIp = result["remote_addr"]
	httpRequest.timestamp = convertDateStringToTime(result["time_local"])
	httpRequest.requestMethod = result["method"]
	httpRequest.requestPath = result["path"]
	httpRequest.httpVersion = result["protocol"]
	httpRequest.httpStatus, _ = strconv.Atoi(result["status"])
	httpRequest.bodyBytesSent, _ = strconv.Atoi(result["body_bytes_sent"])
	httpRequest.httpReferer = result["referer"]
	httpRequest.userAgent = result["user_agent"]
	httpRequest.latency = 0.0
	httpRequest.xForwardedFor = result["http_x_forwarded_for"]

	parseUserAgentAndSetFields(logParser.userAgentParser, httpRequest.userAgent, &httpRequest)
	lookupIpAndSetFields(logParser.ipLookupService, httpRequest.sourceIp, &httpRequest)

	return httpRequest, nil
}
