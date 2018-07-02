package istioutils

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
)

// normalRetryDelay is a desired delay among HTTP requests to the same API
const normalRetryDelay = 3 * time.Second

// minRetryDelay is a min delay among HTTP requests to the same API
const minRetryDelay = 1 * time.Second

const istioPilotV1Registration = "http://istio-pilot.istio-system.svc.cluster.local:8080/v1/registration"

// WaitForIstioSidecar awaits starting up Istio sidecar proxy (a.k.a istio-proxy) in the same Pod (or in the local host)
// until time's up.
func WaitForIstioSidecar(timeout time.Duration) ([]byte, error) {
	return WaitForIstioPilot(timeout)
}

func getPilotRegistration(timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest("GET", istioPilotV1Registration, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	client := http.DefaultClient
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}
	return bytes, nil
}

// WaitForIstioPilot awaits reaching to Istio Pilot (istio-pilot) in another pod until time's up
func WaitForIstioPilot(timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	var err error
	var svcs []byte
	for {
		now := time.Now()
		if now.Before(deadline) {
			svcs, err = getPilotRegistration(deadline.Sub(now))
		} else { // time's up
			if err == nil {
				return nil, fmt.Errorf("too short timeout")
			}
			return nil, err
		}

		// delay the next try
		if err != nil {
			remain := deadline.Sub(time.Now())
			if remain > normalRetryDelay {
				time.Sleep(normalRetryDelay)
			} else if remain > minRetryDelay {
				time.Sleep(minRetryDelay)
			}
		} else {
			return svcs, nil // success
		}
	}
}

func KillEnvoy(envoyAddress string) error {
	const urlKill = `http://%s/quitquitquit`
	res, err := http.Get(fmt.Sprintf(urlKill, envoyAddress))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%s", res.Status)
	}
	return nil
}
