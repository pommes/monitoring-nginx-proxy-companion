package useragent

import (
	"github.com/mssola/useragent"
)

type MssolaParser struct {
}

func (MssolaParser) Parse(userAgentString string) UserAgent {
	ua := useragent.New(userAgentString)

	mobile := ua.Mobile()
	browser, browserVersion := ua.Browser()
	os := ua.OS()
	bot := ua.Bot()

	userAgent := UserAgent{
		Mobile:         mobile,
		Bot:            bot,
		Browser:        browser,
		BrowserVersion: browserVersion,
		Os:             os,
	}
	/*
		if os == "" {
			log.Println("Did not find OS in User Agent String, ", userAgentString)
		}
	*/

	return userAgent
}
