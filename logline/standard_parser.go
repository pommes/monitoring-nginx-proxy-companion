package logline

import (
	"errors"
	"log"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/useragent"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StandardParser struct {
	userAgentParser useragent.Parser
	ipLookupService geoip.IPLocator
}

func (parserarser StandardParser) Parse(logLine string) (HttpRequest, error) {
	const Regex = `^\s*(\S+)\s+(\S+).+\[(.+)\]\s+"([^"]+)"\s+(\S+)\s+(\S+)\s+"([^"]+)"\s+"([^"]+)"($|\s+"([^"]+)"|\s+([0-9.]+))`

	var logLineParserRegex = regexp.MustCompile(Regex)

	logLineParserRegexResult := logLineParserRegex.FindStringSubmatch(logLine)
	if len(logLineParserRegexResult) <= 0 {
		return HttpRequest{}, errors.New("Log line did not match nginx log line.")
	}

	regexFieldIndex := 1
	host := logLineParserRegexResult[regexFieldIndex]
	regexFieldIndex++
	remoteAddress := logLineParserRegexResult[regexFieldIndex]
	regexFieldIndex++
	timestamp := logLineParserRegexResult[regexFieldIndex]

	regexFieldIndex++
	httpRequestHeader := logLineParserRegexResult[regexFieldIndex]
	var httpRequestRegex = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)`)
	httpRequestRegexResult := httpRequestRegex.FindStringSubmatch(httpRequestHeader)
	requestType := httpRequestRegexResult[1]
	requestPath := httpRequestRegexResult[2]
	httpVersion := httpRequestRegexResult[3]

	regexFieldIndex++
	httpStatus, _ := strconv.Atoi(logLineParserRegexResult[regexFieldIndex])
	regexFieldIndex++
	bodyBytesSent, _ := strconv.Atoi(logLineParserRegexResult[regexFieldIndex])
	regexFieldIndex++
	httpReferer := logLineParserRegexResult[regexFieldIndex]
	regexFieldIndex++
	userAgent := logLineParserRegexResult[regexFieldIndex]

	regexFieldIndex += 2
	lastField := logLineParserRegexResult[regexFieldIndex]

	var latencyFloat float64
	var xForwardedFor string
	if lastField != "" {
		var err error
		latencyFloat, err = strconv.ParseFloat(lastField, 64)
		xForwardedFor = ""
		if err != nil {
			// Non-parsable as float because it contains a non-number
			//log.Println("ERROR: reason:", err)
			latencyFloat = 0
			xForwardedFor = strings.Replace(lastField, "\\\"", "", -1)
		}
	}

	httpRequest := HttpRequest{}
	httpRequest.Host = host
	httpRequest.SourceIp = remoteAddress
	httpRequest.Timestamp = convertDateStringToTime(timestamp)
	httpRequest.RequestMethod = requestType
	httpRequest.RequestPath = requestPath
	httpRequest.HttpVersion = httpVersion
	httpRequest.HttpStatus = httpStatus
	httpRequest.BodyBytesSent = bodyBytesSent
	httpRequest.HttpReferer = httpReferer
	httpRequest.UserAgent = userAgent
	httpRequest.Latency = latencyFloat
	httpRequest.xForwardedFor = xForwardedFor

	parseUserAgentAndSetFields(parserarser.userAgentParser, userAgent, &httpRequest)
	lookupIpAndSetFields(parserarser.ipLookupService, remoteAddress, &httpRequest)

	return httpRequest, nil
}

func parseUserAgentAndSetFields(userAgentParser useragent.Parser, userAgentString string, httpRequest *HttpRequest) {
	userAgent := userAgentParser.Parse(userAgentString)

	httpRequest.Browser = userAgent.Browser
	httpRequest.BrowserVersion = userAgent.BrowserVersion
	httpRequest.Os = userAgent.Os
	httpRequest.Mobile = userAgent.Mobile
	httpRequest.Bot = userAgent.Bot
}

func lookupIpAndSetFields(ipLookupService geoip.IPLocator, ip string, httpRequest *HttpRequest) {
	ipLocation := ipLookupService.Lookup(ip)

	httpRequest.Country = ipLocation.Country
	httpRequest.City = ipLocation.City
}

func convertDateStringToTime(dateString string) time.Time {
	const DateFormat = "02/Jan/2006:15:04:05 +0000"
	t, err := time.Parse(DateFormat, dateString)

	if err != nil {
		log.Fatal("Could not parse date string, reason: ", err)
	}

	return t
}
