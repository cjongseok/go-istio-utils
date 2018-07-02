package istioutils

import (
	"fmt"
	"testing"
	"time"
)

func TestWaitForIstioPilot(t *testing.T) {
	timeout := 25 * time.Second
	start := time.Now()
	svcs, err := WaitForIstioPilot(timeout)
	done := time.Now()
	fmt.Println("elapsed:", done.Sub(start))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
		return
	}
	fmt.Println("svcs:", string(svcs))
}
func TestWaitForIstioSidecar(t *testing.T) {
	timeout := 25 * time.Second
	start := time.Now()
	svcs, err := WaitForIstioSidecar(timeout)
	done := time.Now()
	fmt.Println("elapsed:", done.Sub(start))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
		return
	}
	fmt.Println("svcs:", string(svcs))
}
