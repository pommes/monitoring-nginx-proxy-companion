package main

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"nginx-proxy-metrics/config"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/logline"
	"nginx-proxy-metrics/logminer"
	"nginx-proxy-metrics/persistence"
	"nginx-proxy-metrics/useragent"
)

func main() {
	log.Println("Loading config.")
	config.LoadConfig()

	log.Println("Creating dependencies of log miner.")

	log.Println("Creating persistence client.")
	httpRequestPersister := persistence.InfluxdbHttpRequestPersister{}
	dbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: config.InfluxUrl,
	})

	if err != nil {
		log.Fatal("Could not setup persistence client, reason: ", err)
	}
	httpRequestPersister.NewInfluxdbHttpRequestPersister(dbClient)

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
