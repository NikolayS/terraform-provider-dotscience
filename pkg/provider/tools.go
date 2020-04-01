package provider

import (
	"context"
	"fmt"
	"time"
)

func waitOnFunction(name string, timeout, interval time.Duration, testingFunction func() bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Took more than %d seconds for %s to be ready.", timeout/time.Second, name)
		default:
			isReady := testingFunction()
			if isReady {
				return nil
			}
			time.Sleep(interval)
		}
	}
}
