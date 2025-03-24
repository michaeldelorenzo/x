// Package testutils provides helpful utilities for writing tests and their assertions
package testutils

import (
	"github.com/stretchr/testify/require"
)

// RequireBothNilOrEqual requires that either both expected and actual are nil or that *input and *output are equal via
// require.Equal
func RequireBothNilOrEqual[T any](t require.TestingT, expected, actual *T) {
	RequireBothNilOrEqualFn(t, expected, actual, func(t require.TestingT, expected, actual T) {
		require.Equal(t, expected, actual)
	})
}

// RequireBothNilOrEqualFn requires that either both expected and actual are nil or that *input and *output are equal via
// a custom equality callback function
func RequireBothNilOrEqualFn[T any](t require.TestingT, expected, actual *T, requireEqualFn func(t require.TestingT, expected, actual T)) {
	if expected != nil {
		require.NotNil(t, actual)
		requireEqualFn(t, *expected, *actual)
	} else {
		require.Nil(t, actual)
	}
}
