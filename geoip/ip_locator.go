package geoip

type IPLocation struct {
	Country string
	City    string
}

type IPLocator interface {
	Lookup(ip string) IPLocation
}
