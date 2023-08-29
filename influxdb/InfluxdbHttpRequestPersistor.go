package influxdb

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"nginx-proxy-metrics/config"
	"nginx-proxy-metrics/logparser"
	"strconv"
)

const SERIES_NAME = "http_requests"

type InfluxdbHttpRequestPersistor struct {
	influxClient client.Client
}

func (influxdbLogPersistor *InfluxdbHttpRequestPersistor) Setup() {
	dbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: config.InfluxUrl,
	})

	if err != nil {
		log.Fatal("Could not setup influxdb client, reason: ", err)
	}

	_, db_err := queryDB(dbClient, fmt.Sprintf("CREATE DATABASE %s", config.InfluxDbName))
	if db_err != nil {
		log.Fatal("Could not create database, reason: ", db_err)
	}

	log.Println(fmt.Sprintf("  - Altering retention policy %s_retention to %s", config.InfluxDbName, config.InfluxDbRetentionDuration))
	_, db_err = queryDB(dbClient, fmt.Sprintf("ALTER RETENTION POLICY %s_retention ON %s DURATION %s DEFAULT",
		config.InfluxDbName, config.InfluxDbName, config.InfluxDbRetentionDuration))
	if db_err != nil {
		log.Println("  - Could not ALTER retention policy, reason: ", db_err)

		log.Println(fmt.Sprintf("  - Creating retention policy %s_retention with %s", config.InfluxDbName, config.InfluxDbRetentionDuration))
		_, db_err = queryDB(dbClient, fmt.Sprintf("CREATE RETENTION POLICY %s_retention ON %s DURATION %s REPLICATION 1 DEFAULT",
			config.InfluxDbName, config.InfluxDbName, config.InfluxDbRetentionDuration))
		if db_err != nil {
			log.Println("  - Could not CREATE retention policy, reason: ", db_err)
		}
	}

	influxdbLogPersistor.influxClient = dbClient
}

func (influxLogPersistor InfluxdbHttpRequestPersistor) Persist(httpRequest logparser.HttpRequest) {
	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: config.InfluxDbName,
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{
		"host":            httpRequest.Host,
		"request_method":  httpRequest.RequestMethod,
		"http_version":    httpRequest.HttpVersion,
		"http_status":     strconv.Itoa(httpRequest.HttpStatus),
		"browser":         httpRequest.Browser,
		"browser_version": httpRequest.BrowserVersion,
		"os":              httpRequest.Os,
		"mobile":          strconv.FormatBool(httpRequest.Mobile),
		"country":         httpRequest.Country,
		"city":            httpRequest.City,
		"instance":        config.InfluxDbTagInstance,
		"bot":             strconv.FormatBool(httpRequest.Bot),
	}

	fields := map[string]interface{}{
		"source_ip":       httpRequest.SourceIp,
		"request_path":    httpRequest.RequestPath,
		"body_bytes_sent": httpRequest.BodyBytesSent,
		"http_referer":    httpRequest.HttpReferer,
		"user_agent":      httpRequest.UserAgent,
		"latency":         httpRequest.Latency,
	}

	point, err := client.NewPoint(SERIES_NAME, tags, fields, httpRequest.Timestamp)
	if err != nil {
		log.Fatal(err)
	}

	batchPoints.AddPoint(point)

	if err := influxLogPersistor.influxClient.Write(batchPoints); err != nil {
		log.Println("Could not insert into influxdb, reason:", err)
		return
	}
}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command: cmd,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
