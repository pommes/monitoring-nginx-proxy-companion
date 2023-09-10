package config

import (
	"fmt"
	"log"
	"os"
)

const (
	EnvProxyContainerName        = "PROXY_CONTAINER_NAME"
	EnvInfluxUrl                 = "INFLUX_URL"
	EnvInfluxDbName              = "INFLUX_DB_NAME"
	EnvInfluxDbRetentionDuration = "INFLUX_DB_RETENTION_DURATION"
	EnvInfluxDbTagInstance       = "INFLUX_DB_TAG_INSTANCE"
	EnvInfluxDbTagSourceIpsLocal = "INFLUX_DB_TAG_SOURCE_IPS_LOCAL"
)

var (
	ProxyContainerName        string
	InfluxUrl                 string
	InfluxDbName              string
	InfluxDbRetentionDuration string
	InfluxDbTagInstance       string
	InfluxDbTagSourceIpsLocal string
)

func LoadConfig() {
	ProxyContainerName = getEnvOrPanic(EnvProxyContainerName)
	InfluxUrl = getEnvOrPanic(EnvInfluxUrl)
	InfluxDbName = getEnvOrPanic(EnvInfluxDbName)
	InfluxDbRetentionDuration = getEnvOrPanic(EnvInfluxDbRetentionDuration)
	InfluxDbTagInstance = getEnvOrPanic(EnvInfluxDbTagInstance)
	InfluxDbTagSourceIpsLocal = getEnvOrPanic(EnvInfluxDbTagSourceIpsLocal)
}

func getEnvOrPanic(envName string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		panic(fmt.Sprintf("ERROR: Environment variable '%s' not set", envName))
	}
	log.Printf("Init config '%s' = '%s'\n", envName, envValue)

	return envValue
}
