package main

import (
	"log"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/logline"
	"nginx-proxy-metrics/logminer"
	"nginx-proxy-metrics/persistence"
	"nginx-proxy-metrics/useragent"
)

func main() {
	log.Println("Starting nginx-proxy-metrics.")

	log.Println("Creating dependencies of log miner.")

	log.Println("Creating persistence client.")
	httpRequestPersister := persistence.InfluxdbHttpRequestPersister{}
	httpRequestPersister.Setup()

	log.Println("Creating user agent parser, ip lookup service and log parser.")
	loglineParser := logline.ProxyParser{
		UserAgentParser: useragent.MssolaParser{},
		IPLocator:       geoip.GeoLite2Locator{},
	}

	log.Println("Creating docker container log miner.")
	logMiner := logminer.DockerContainerLogMiner{
		LoglineParser:        loglineParser,
		HttpRequestPersister: httpRequestPersister,
	}

	log.Println("Start log mining.")
	logMiner.Mine()
}
