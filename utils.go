package istioutils

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"io"
)

// maxLocalHttpTimeout is the max timeout for HTTP request inside pod
const maxLocalHttpTimeout = 3 * time.Second

// maxClusterHttpTimeout is the max timeout of HTTP request inside cluster
const maxCluterHttpTimeout = 7 * time.Second

// normalRetryDelay is a desired delay among HTTP requests to the same API
const normalRetryDelay = 3 * time.Second

// minRetryDelay is a min delay among HTTP requests to the same API
const minRetryDelay = 1 * time.Second

// Istio URLs
const (
	localEnvoyURL          = "http://localhost:15000"
	killLocalEnvoyURL      = localEnvoyURL + "/quitquitquit"
	pilotURL               = "http://istio-pilot.istio-system.svc.cluster.local:8080"
	pilotV1RegistrationURL = pilotURL + "/v1/registration"
)

var logWriter io.Writer

// SetLogging sets up log writer
func SetLogging(writer io.Writer) {
	logWriter = writer
}

// WaitForSidecarProxy awaits starting up Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host)
// until time's up.
func WaitForSidecarProxy(timeout time.Duration) error {
	return httpGetUntil200(localEnvoyURL, timeout, maxLocalHttpTimeout)
}

// WaitForPilot awaits reaching to Istio Pilot (istio-pilot) in another pod until time's up
func WaitForPilot(timeout time.Duration) error {
	return httpGetUntil200(pilotV1RegistrationURL, timeout, maxCluterHttpTimeout)
}

// KillSidecarProxy kills Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host)
func KillSidecarProxy(timeout time.Duration) error {
	return httpGetUntil200(killLocalEnvoyURL, timeout, maxLocalHttpTimeout)
}

func httpGetWithTimeout(url string, timeout time.Duration) (res *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	client := http.DefaultClient
	res, err = client.Do(req.WithContext(ctx))
	return
}

func minDuration(a, b time.Duration) time.Duration {
	if a > b {
		return b
	}
	return a
}

func httpGetUntil200(url string, timeout time.Duration, maxHttpTimeout time.Duration) (err error) {
	var res *http.Response
	deadline := time.Now().Add(timeout)
	fastFail := false
	for {
		now := time.Now()
		if !fastFail && now.Before(deadline) {
			res, err = httpGetWithTimeout(url, minDuration(deadline.Sub(now), maxHttpTimeout))
		} else { // time's up
			if err == nil {
				if res == nil {
					err = fmt.Errorf("timeout")
				} else {
					err = fmt.Errorf("timeout; last status code = %d", res.StatusCode)
				}
			}
			return
		}

		// delay the next try
		if err == nil && res != nil && res.StatusCode == http.StatusOK {
			return
		}
		if logWriter != nil {
			fmt.Fprintf(logWriter, "http Get %s failure", url)
			if res != nil {
				fmt.Fprintf(logWriter, "; %s", res.Status)
			}
			if err != nil{
				fmt.Fprintf(logWriter, "; %v", err)
			}
			fmt.Fprint(logWriter, "\n")
		}
		remain := deadline.Sub(time.Now())
		if remain > normalRetryDelay {
			time.Sleep(normalRetryDelay)
		} else if remain > minRetryDelay {
			time.Sleep(minRetryDelay)
		} else {
			fastFail = true
		}
	}
}
