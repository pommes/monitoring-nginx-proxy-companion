package config

import (
	"fmt"
	"log"
	"os"
)

var (
	ProxyContainerName        = getEnvOrPanic("PROXY_CONTAINER_NAME")
	InfluxUrl                 = getEnvOrPanic("INFLUX_URL")
	InfluxDbName              = getEnvOrPanic("INFLUX_DB_NAME")
	InfluxDbRetentionDuration = getEnvOrPanic("INFLUX_DB_RETENTION_DURATION")
	InfluxDbTagInstance       = getEnvOrPanic("INFLUX_DB_TAG_INSTANCE")
	LocalSourceIPs            = getEnvOrPanic("INFLUX_DB_TAG_SOURCE_IPS_LOCAL")
)

func getEnvOrPanic(envName string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		panic(fmt.Sprintf("ERROR: Environment variable '%s' not set", envName))
	}
	log.Printf("Init config '%s' = '%s'\n", envName, envValue)

	return envValue
}
