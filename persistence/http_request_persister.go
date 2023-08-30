package persistence

import (
	"nginx-proxy-metrics/logline"
)

type HttpRequestPersister interface {
	Persist(httpRequest logline.HttpRequest)
}
