package testutils

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// This file contains testhelpers that allow you to test other test utilities

// MockTestingT is a wrapper on testing.T that overrides `FailNow` to simply set a flag as opposed to
// actually rejecting the test. Useful for testing test utilities.
type MockTestingT struct {
	failNowCalled bool
	T             *testing.T
	require.TestingT
}

func (m *MockTestingT) FailNow() {
	// register the method is called
	m.failNowCalled = true
	// exit, as normal behaviour
	runtime.Goexit()
}

func (m *MockTestingT) Errorf(format string, args ...interface{}) {
	// no-op
}

// FailNowCalled return true if FailNow was called
func (m *MockTestingT) FailNowCalled() bool {
	return m.failNowCalled
}

func (m *MockTestingT) Reset() {
	m.failNowCalled = false
}

// RequireTestFailure verifies that a test fails and should wrap the callback passed to a t.Run call
func RequireTestFailure(testFn func(t require.TestingT)) func(*testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup

		// create a mock structure for TestingT
		mockT := &MockTestingT{T: t}

		// setup the barrier
		wg.Add(1)
		// start a co-routine to execute the test function f
		// and release the barrier at its end
		go func() {
			defer wg.Done()
			testFn(mockT)
		}()

		// wait for the barrier.
		wg.Wait()
		// verify fail now is invoked
		require.True(t, mockT.FailNowCalled())
	}

}

// RequireNoTestFailure verifies that a test passes and should wrap the callback passed to a t.Run call
func RequireNoTestFailure(testFn func(t require.TestingT)) func(*testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup

		// create a mock structure for TestingT
		mockT := &MockTestingT{T: t}

		// setup the barrier
		wg.Add(1)
		// start a co-routine to execute the test function f
		// and release the barrier at its end
		go func() {
			defer wg.Done()
			testFn(mockT)
		}()

		// wait for the barrier.
		wg.Wait()
		// verify fail now is invoked
		require.False(t, mockT.FailNowCalled())
	}
}
