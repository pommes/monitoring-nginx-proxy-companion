package useragent

type UserAgent struct {
	Browser        string
	BrowserVersion string
	Os             string
	Mobile         bool
	Bot            bool
}

type Parser interface {
	Parse(userAgent string) UserAgent
}
