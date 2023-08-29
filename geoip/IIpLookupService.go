package geoip

type IPLocation struct {
	Country string
	City    string
}

type IIpLookupService interface {
	Lookup(ip string) IPLocation
}
