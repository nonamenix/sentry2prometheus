package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

func init() {
	prometheus.MustRegister(version.NewCollector("sentry2prometheus"))
}

func main() {
	var (
		sentryURL          = flag.String("sentry-url", "https://sentry.io", "The sentry url")
		organization       = flag.String("organization", "wargaming", "Organization name in sentry")
		statsPeriod        = flag.String("stats-period", "24h", "Sentry stats period")
		query              = flag.String("query", "", "Sentry query for projects filtering")
		listenAddress      = flag.String("port", ":9412", "The address to listen on for HTTP requests.")
		authorizationToken = flag.String("token", "", "Sentry API authorization token")
	)
	flag.Parse()

	var config = Config{
		sentryURL:          *sentryURL,
		organization:       *organization,
		query:              *query,
		statsPeriod:        *statsPeriod,
		authorizationToken: *authorizationToken,
	}

	log.Infoln("Starting sentry2prometheus", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/",
		func(responseWriter http.ResponseWriter, request *http.Request) {
			SentryOrganizationMetricsHandler(config, responseWriter)
		})

	log.Infoln("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func SentryOrganizationMetricsHandler(config Config, responseWriter http.ResponseWriter) {
	start := time.Now()

	projects, err := fetchErrorsFromSentryHandler(config)
	if err != nil {
		fmt.Fprintln(responseWriter, "probe_success 0")
		log.Errorf("Error on sentry projects fetching %s", err)
		return
	}

	fmt.Fprintf(responseWriter, "# HELP probe_sentry_errors_received Errors count since timestamp\n")
	fmt.Fprintf(responseWriter, "# TYPE probe_sentry_errors_received counter\n")

	for _, project := range projects {
		stat := project.Stats[len(project.Stats)-1]
		timestamp := int(stat[0])
		errorsCount := int(stat[1])

		fmt.Fprintf(responseWriter, "probe_sentry_errors_received{project=\"%s\", timestamp=\"%d\"} %d\n", project.Slug, timestamp, errorsCount)
	}
	fmt.Fprintln(responseWriter, "probe_success 1")
	fmt.Fprintf(responseWriter, "probe_projects_count %d\n", len(projects))
	fmt.Fprintf(responseWriter, "probe_duration_seconds %f\n", time.Since(start).Seconds())

	return
}

type Project struct {
	Slug  string      `json:"slug"`
	Stats [][]float64 `json:"stats"`
	Team  struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	} `json:"team"`
}

type Config struct {
	sentryURL          string
	organization       string
	statsPeriod        string
	query              string
	authorizationToken string
	timeout            time.Duration
}

func fetchErrorsFromSentryHandler(config Config) ([]Project, error) {
	requestURL := config.sentryURL + "/api/0/organizations/" + config.organization + "/projects/?statsPeriod=24h&query=" + config.query
	log.Infof("requestURL %s", requestURL)

	client := &http.Client{}

	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Errorf("Error creating request for target %s", err)
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+config.authorizationToken)

	resp, err := client.Do(request)
	if err != nil && resp == nil {
		log.Warnf("Error for HTTP request to %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var projects = []Project{}

	decodeErr := json.NewDecoder(resp.Body).Decode(&projects)
	if decodeErr != nil {
		log.Errorf("Error on json parsing %s", err)
		return nil, decodeErr
	}

	return projects, nil

}
