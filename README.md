go-istio-utils
===
Istio utils in Go

Utils
---
### WaitForIstioSidecar(timeout)
It waits for starting up Istio sidecar proxy (a.k.a istio-proxy) in the same Pod (or in the local host) until time's up.

### WaitForIstioPilot(timeout)
It waits for reaching to Istio Pilot (istio-pilot) in another pod until time's up.
