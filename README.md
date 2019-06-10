# Sentry Errors Prometheus Exporter 

```bash
./sentry2prometheus --help
Usage of ./sentry2prometheus:
  -organization string
    	Organization name in sentry (default "XXX")
  -port string
    	The address to listen on for HTTP requests. (default ":9412")
  -query string
    	Sentry query for projects filtering
  -sentry-url string
    	The sentry url (default "https://sentry.io")
  -stats-period string
    	Sentry stats period (default "24h")
  -token string
    	Sentry API authorization token
```

## Build

```bash
go install
go build
```

## Usage Example

```bash
./sentry2prometheus 
    --sentry-url=https://sentry.io 
    --organization=XXX 
    --query=team:project 
    --token=7daef5d63f6746ae8b1f5abe2e3872786ee7cea23ade46e29b536c28463ebe
```

Visiting [http://localhost:9412/](http://localhost:9412/) will return metrics for a the sentry projects in your `organization` filtered by `query`
