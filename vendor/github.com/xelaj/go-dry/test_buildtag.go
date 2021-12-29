// +build test

package dry

var _ = func() bool {
	testMode = true
	return true
}()
