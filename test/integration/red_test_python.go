//go:build integration

package integration

import (
	"net"
	"testing"

	"github.com/mariomac/guara/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/beyla/test/integration/components/prom"
)

func testREDMetricsForPythonHTTPLibrary(t *testing.T, url string, comm string) {
	urlPath := "/greeting"

	// Call 3 times the instrumented service, forcing it to:
	// - take a large JSON file
	// - returning a 200 code
	for i := 0; i < 4; i++ {
		doHTTPGet(t, url+urlPath, 200)
	}

	// Eventually, Prometheus would make this query visible
	pq := prom.Client{HostPort: prometheusHostPort}
	var results []prom.Result
	test.Eventually(t, testTimeout, func(t require.TestingT) {
		var err error
		results, err = pq.Query(`http_server_duration_seconds_count{` +
			`http_request_method="GET",` +
			`http_response_status_code="200",` +
			`service_namespace="integration-test",` +
			`service_name="` + comm + `",` +
			`url_path="` + urlPath + `"}`)
		require.NoError(t, err)
		enoughPromResults(t, results)
		val := totalPromCount(t, results)
		assert.LessOrEqual(t, 3, val)
		if len(results) > 0 {
			res := results[0]
			addr := net.ParseIP(res.Metric["client_address"])
			assert.NotNil(t, addr)
		}
	})
}

func testREDMetricsPythonHTTP(t *testing.T) {
	for _, testCaseURL := range []string{
		"http://localhost:8081",
	} {
		t.Run(testCaseURL, func(t *testing.T) {
			waitForTestComponents(t, testCaseURL)
			testREDMetricsForPythonHTTPLibrary(t, testCaseURL, "python3.11") // reusing what we do for NodeJS
		})
	}
}

func testREDMetricsPythonHTTPS(t *testing.T) {
	for _, testCaseURL := range []string{
		"https://localhost:8081",
	} {
		t.Run(testCaseURL, func(t *testing.T) {
			waitForTestComponents(t, testCaseURL)
			testREDMetricsForPythonHTTPLibrary(t, testCaseURL, "python3.11") // reusing what we do for NodeJS
		})
	}
}
