package logparser

import (
	"errors"
	"log"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/useragentparser"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StandardLogParser struct {
	userAgentParser useragentparser.IUserAgentParser
	ipLookupService geoip.IIpLookupService
}

const LOG_LINE_REGEX = `^\s*(\S+)\s+(\S+).+\[(.+)\]\s+"([^"]+)"\s+(\S+)\s+(\S+)\s+"([^"]+)"\s+"([^"]+)"($|\s+"([^"]+)"|\s+([0-9.]+))`

const STANDARD_LOG_LINE_DATE_FORMAT = "02/Jan/2006:15:04:05 +0000"

func (standardLogParser StandardLogParser) Parse(logLine string) (HttpRequest, error) {
	var logLineParserRegex = regexp.MustCompile(LOG_LINE_REGEX)

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

	parseUserAgentAndSetFields(standardLogParser.userAgentParser, userAgent, &httpRequest)
	lookupIpAndSetFields(standardLogParser.ipLookupService, remoteAddress, &httpRequest)

	return httpRequest, nil
}

func parseUserAgentAndSetFields(userAgentParser useragentparser.IUserAgentParser, userAgentString string, httpRequest *HttpRequest) {
	userAgent := userAgentParser.Parse(userAgentString)

	httpRequest.Browser = userAgent.Browser
	httpRequest.BrowserVersion = userAgent.BrowserVersion
	httpRequest.Os = userAgent.Os
	httpRequest.Mobile = userAgent.Mobile
	httpRequest.Bot = userAgent.Bot
}

func lookupIpAndSetFields(ipLookupService geoip.IIpLookupService, ip string, httpRequest *HttpRequest) {
	ipLocation := ipLookupService.Lookup(ip)

	httpRequest.Country = ipLocation.Country
	httpRequest.City = ipLocation.City
}

func convertDateStringToTime(dateString string) time.Time {
	t, err := time.Parse(STANDARD_LOG_LINE_DATE_FORMAT, dateString)

	if err != nil {
		log.Fatal("Could not parse date string, reason: ", err)
	}

	return t
}
