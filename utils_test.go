package istioutils

import (
	"fmt"
	"testing"
	"time"
	"os"
)

func TestWaitForPilot(t *testing.T) {
	timeout := 25 * time.Second
	start := time.Now()
	err := WaitForPilot(timeout)
	done := time.Now()
	fmt.Println("elapsed:", done.Sub(start))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
		return
	}
}

func TestWaitForSidecarProxy(t *testing.T) {
	SetLogging(os.Stdout)
	timeout := 25 * time.Second
	start := time.Now()
	err := WaitForSidecarProxy(timeout)
	done := time.Now()
	fmt.Println("elapsed:", done.Sub(start))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
		return
	}
}
