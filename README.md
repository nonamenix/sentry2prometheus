# Sentry Errors Prometheus Exporter 

```bash
./sentry2prometheus --help

Usage of ./sentry2prometheus:
      --extra-labels strings   Extra labels for prometheus metrics splitted by ':'
      --organization string    Organization name in sentry (default "XXX")
      --port string            The address to listen on for HTTP requests. (default ":9412")
      --query string           Sentry query for projects filtering
      --sentry-url string      The sentry url (default "https://sentry.io")
      --stats-period string    Sentry stats period (default "24h")
      --token string           Sentry API authorization token

```

## Build

```bash
go install
go build
```

## Usage Example

Take sentry token from https://sentry.io/api

```bash
./sentry2prometheus --sentry-url=https://sentry.io \
    --organization=XXX \
    --query=team:web \
    --token=token_from_sentry \
    --extra-labels=team:web
```

Visiting [http://localhost:9412/](http://localhost:9412/) will return metrics for a the sentry projects in your `organization` filtered by `query`

```text
# HELP probe_sentry_errors_received Errors count since timestamp
# TYPE probe_sentry_errors_received counter
probe_sentry_errors_received{project="portal", timestamp=1560322800, team="web"} 5
probe_success 1
probe_projects_count 1
probe_duration_seconds 0.132
```
