FROM alpine:latest

LABEL org.opencontainers.image.description="nginx-proxy-metrics parses nginx-proxy access.logs from the proxy containers stdout, adds location info to the data and stores them in an influxDB that can be used as a data source for grafana."

ENV PROXY_CONTAINER_NAME nginx
ENV INFLUX_DB_NAME monitoring
ENV INFLUX_DB_RETENTION_DURATION 8w
ENV INFLUX_DB_TAG_SOURCE_IPS_LOCAL "127., 10., 192., 176."

VOLUME /GeoLite2-City.mmdb
COPY .build/main /main

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]