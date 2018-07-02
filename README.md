go-istio-utils
===
Istio utils in Go

Utils
---
### WaitForSidecarProxy(timeout)
It waits for starting up Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host) until time's up.

### WaitForPilot(timeout)
It waits for reaching to Istio Pilot (istio-pilot) in another pod until time's up.

### KillSidecarProxy(timeout)
It kills Istio sidecar proxy (a.k.a envoy or istio-proxy) in the same Pod (or in the local host).
It gives up the killing when time's up.
