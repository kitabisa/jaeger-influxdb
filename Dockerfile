FROM jaegertracing/all-in-one:1.22.0

ENV COLLECTOR_ZIPKIN_HTTP_PORT=9441
ENV SPAN_STORAGE_TYPE="grpc-plugin"
ENV GRPC_STORAGE_PLUGIN_BINARY="/opt/jaeger-influxdb/jaeger-influxdb-linux"

COPY ./bin/jaeger-influxdb/jaeger-influxdb-linux /opt/jaeger-influxdb/jaeger-influxdb-linux

ENTRYPOINT ["/go/bin/all-in-one-linux"]
CMD ["--sampling.strategies-file=/etc/jaeger/sampling_strategies.json"]
