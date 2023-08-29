package useragentparser

type UserAgent struct {
	Browser        string
	BrowserVersion string
	Os             string
	Mobile         bool
	Bot            bool
}

type IUserAgentParser interface {
	Parse(userAgent string) UserAgent
}
