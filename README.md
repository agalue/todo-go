# TODO API with PostgreSQL Backend

This repository contains a simple Go application of a TODO application with a ReST API using pure `net/http` using the new routing capabilities available in version 1.22, a PostgreSQL-backed layer based on [GORM](https://gorm.io/), and a simple WebUI implemented in [VueJS](https://vuejs.org/) with [Vuetify](https://vuetifyjs.com/), all embedded in a single binary with a size-optimized Docker image.

It also has a Middleware layer that exposes basic metrics for Prometheus and traces using OpenTelemetry gRPC API and Swagger docs generated externally.

For end-to-end testing, a Docker Compose compiles and deploys the application and all the dependencies (PostgreSQL, Prometheus, Grafana Tempo, Grafana).

## Update Swagger Docs

From the project's folder, run the following

```bash
docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest init -g app/handlers.go -o app/docs
```

## Development Environment

We already have a Docker Compose file for the end-to-end testing, and we can use that as our development environment via [Tilt](https://tilt.dev/).

Once you have the CLI installed (as well as Docker), you can start it with:

```bash
tilt up
```

Any time you change something on the source code, a new container image will be built automatically with your change (you can track progress through Tilt's UI and even check the logs of each container).

To tear down:

```bash
tilt down
```

## Monitoring using RED Method

**Rate:**

```promql
sum(
  rate(
    http_requests_total{}[1m]
  )
)
```

or

```promql
sum(
  rate(
    http_request_duration_seconds_count{}[1m]
  )
)
```

**Errors:**

```promql
sum(
  rate(
    http_requests_total{status !~ "2.."}[1m]
  )
)
```

or

```promql
sum(
  rate(
    http_request_duration_seconds_count{status !~ "2.."}[1m]
  )
)
```

**Duration:**

```promql
histogram_quantile(0.99,
  sum(
    rate(
      http_request_duration_seconds_bucket{}[1m]
    )
  ) by (le)
)
```

## Tracing

Use the following TraceQL query to look at the application traces:

```traceql
{.http.target=~"/api/v1/.*"}
```
