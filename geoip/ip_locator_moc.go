package geoip

type IPLocatorMock struct{}

func (mock IPLocatorMock) Lookup(ip string) IPLocation {
	return IPLocation{}
}
