package main

import (
	"log"
	"nginx-proxy-metrics/geoip"
	"nginx-proxy-metrics/influxdb"
	"nginx-proxy-metrics/logminer"
	"nginx-proxy-metrics/logparser"
	"nginx-proxy-metrics/useragentparser"
)

func main() {
	log.Println("Starting nginx-proxy-metrics.")

	log.Println("Creating dependencies of log miner.")

	log.Println("Creating influxdb client.")
	influxdbHttpRequestPersistor := influxdb.InfluxdbHttpRequestPersistor{}
	influxdbHttpRequestPersistor.Setup()

	log.Println("Creating user agent parser, ip lookup service and log parser.")
	mssolaUserAgentParser := useragentparser.MssolaUserAgentParser{}
	geoIp2IpLookupService := geoip.GeoIp2IpLookupService{}
	logParser := logparser.ProxyLogParser{
		UserAgentParser: mssolaUserAgentParser,
		IpLookupService: geoIp2IpLookupService,
	}

	log.Println("Creating docker container log miner.")
	dockerContainerLogMiner := logminer.DockerContainerLogMiner{
		LogParser:            logParser,
		HttpRequestPersistor: influxdbHttpRequestPersistor,
	}

	log.Println("Start log mining.")
	dockerContainerLogMiner.Mine()
}
