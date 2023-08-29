package main

import (
	"github.com/mssola/useragent"
)

type MssolaUserAgentParser struct {
}

func (MssolaUserAgentParser) Parse(userAgentString string) UserAgent {
	ua := useragent.New(userAgentString)

	mobile := ua.Mobile()
	browser, browserVersion := ua.Browser()
	os := ua.OS()
	bot := ua.Bot()

	userAgent := UserAgent{
		mobile:         mobile,
		bot:            bot,
		browser:        browser,
		browserVersion: browserVersion,
		os:             os,
	}
	/*
		if os == "" {
			log.Println("Did not find OS in User Agent String, ", userAgentString)
		}
	*/

	return userAgent
}
