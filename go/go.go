package gogo

import "sync"

// SafeGo launches a new goroutine that executes the given function,
// and safely ties it to the provided sync.WaitGroup.
//
// It ensures that the WaitGroup counter is incremented before launching the goroutine,
// and automatically decremented via defer once the function completes execution.
// This pattern prevents common bugs where a missing wg.Done() causes indefinite blocking.
//
// Parameters:
//   - wg: a pointer to a sync.WaitGroup that tracks concurrent tasks.
//   - f: a no-arg function to execute asynchronously.
//
// Example:
//
//	var wg sync.WaitGroup
//	SafeGo(&wg, func() {
//	    // do some work
//	})
//	wg.Wait() // wait for all SafeGo() goroutines to finish
func SafeGo(
	wg *sync.WaitGroup,
	f func(),
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}
