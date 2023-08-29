package main

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
)

const SERIES_NAME = "http_requests"

type InfluxdbHttpRequestPersistor struct {
	influxClient client.Client
}

func (influxdbLogPersistor *InfluxdbHttpRequestPersistor) Setup() {
	log.Println("  - Using INFLUX_URL =", getInfluxUrl())
	log.Println("  - Using INFLUX_DB_NAME =", getInfluxDbName())
	log.Println("  - Using INFLUX_DB_RETENTION_DURATION =", getInfluxDbRetentionDuration())
	log.Println("  - Using INFLUX_DB_TAG_INSTANCE =", getInfluxDbTagInstance())
	dbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: getInfluxUrl(),
	})

	if err != nil {
		log.Fatal("Could not setup influx client, reason: ", err)
	}

	_, db_err := queryDB(dbClient, fmt.Sprintf("CREATE DATABASE %s", getInfluxDbName()))
	if db_err != nil {
		log.Fatal("Could not create database, reason: ", db_err)
	}

	log.Println(fmt.Sprintf("  - Altering retention policy %s_retention to %s", getInfluxDbName(), getInfluxDbRetentionDuration()))
	_, db_err = queryDB(dbClient, fmt.Sprintf("ALTER RETENTION POLICY %s_retention ON %s DURATION %s DEFAULT",
		getInfluxDbName(), getInfluxDbName(), getInfluxDbRetentionDuration()))
	if db_err != nil {
		log.Println("  - Could not ALTER retention policy, reason: ", db_err)

		log.Println(fmt.Sprintf("  - Creating retention policy %s_retention with %s", getInfluxDbName(), getInfluxDbRetentionDuration()))
		_, db_err = queryDB(dbClient, fmt.Sprintf("CREATE RETENTION POLICY %s_retention ON %s DURATION %s REPLICATION 1 DEFAULT",
			getInfluxDbName(), getInfluxDbName(), getInfluxDbRetentionDuration()))
		if db_err != nil {
			log.Println("  - Could not CREATE retention policy, reason: ", db_err)
		}
	}

	influxdbLogPersistor.influxClient = dbClient
}

func (influxLogPersistor InfluxdbHttpRequestPersistor) Persist(httpRequest HttpRequest) {
	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: getInfluxDbName(),
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{
		"host":            httpRequest.host,
		"request_method":  httpRequest.requestMethod,
		"http_version":    httpRequest.httpVersion,
		"http_status":     strconv.Itoa(httpRequest.httpStatus),
		"browser":         httpRequest.browser,
		"browser_version": httpRequest.browserVersion,
		"os":              httpRequest.os,
		"mobile":          strconv.FormatBool(httpRequest.mobile),
		"country":         httpRequest.country,
		"city":            httpRequest.city,
		"instance":        getInfluxDbTagInstance(),
		"bot":             strconv.FormatBool(httpRequest.bot),
	}

	fields := map[string]interface{}{
		"source_ip":       httpRequest.sourceIp,
		"request_path":    httpRequest.requestPath,
		"body_bytes_sent": httpRequest.bodyBytesSent,
		"http_referer":    httpRequest.httpReferer,
		"user_agent":      httpRequest.userAgent,
		"latency":         httpRequest.latency,
	}

	point, err := client.NewPoint(SERIES_NAME, tags, fields, httpRequest.timestamp)
	if err != nil {
		log.Fatal(err)
	}

	batchPoints.AddPoint(point)

	if err := influxLogPersistor.influxClient.Write(batchPoints); err != nil {
		log.Println("Could not insert into influx, reason:", err)
		return
	}
}

func getInfluxUrl() string {
	return GetEnvOrPanic(INFLUX_URL_ENV_NAME)
}

func getInfluxDbName() string {
	return GetEnvOrPanic(INFLUX_DB_ENV_NAME)
}

func getInfluxDbRetentionDuration() string {
	return GetEnvOrPanic(INFLUX_DB_RETENTION_DURATION)
}

func getInfluxDbTagInstance() string {
	return GetEnvOrPanic(INFLUX_DB_TAG_INSTANCE)
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
