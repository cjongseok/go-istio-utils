package istioutils

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

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

// WaitForSidecarProxy awaits starting up Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host)
// until time's up.
func WaitForSidecarProxy(timeout time.Duration) error {
	return httpGetUntil200(localEnvoyURL, timeout)
}

// WaitForPilot awaits reaching to Istio Pilot (istio-pilot) in another pod until time's up
func WaitForPilot(timeout time.Duration) error {
	return httpGetUntil200(pilotV1RegistrationURL, timeout)
}

// KillSidecarProxy kills Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host)
func KillSidecarProxy(timeout time.Duration) error {
	return httpGetUntil200(killLocalEnvoyURL, timeout)
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

func httpGetUntil200(url string, timeout time.Duration) (err error) {
	var res *http.Response
	deadline := time.Now().Add(timeout)
	for {
		now := time.Now()
		if now.Before(deadline) {
			res, err = httpGetWithTimeout(url, timeout)
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
		remain := deadline.Sub(time.Now())
		if remain > normalRetryDelay {
			time.Sleep(normalRetryDelay)
		} else if remain > minRetryDelay {
			time.Sleep(minRetryDelay)
		}
	}
}
