package influxdb

import (
	"nginx-proxy-metrics/logparser"
)

type IHttpRequestPersistor interface {
	Persist(httpRequest logparser.HttpRequest)
}
