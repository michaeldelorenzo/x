// Package testutils provides helpful utilities for writing tests and their assertions
package testutils_test

import (
	"testing"

	"github.com/koneksahealth/x/pkg/testutils"
	"github.com/koneksahealth/x/pkg/utils/ptr"
	"github.com/stretchr/testify/require"
)

func Test_RequireBothNilOrEqual(t *testing.T) {
	t.Run("passes when both pointers are non-nil and equal", testutils.RequireNoTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		actual := ptr.Addr("test")

		// when
		testutils.RequireBothNilOrEqual(t, expected, actual)
	}))

	t.Run("passes when both pointers are nil", testutils.RequireNoTestFailure(func(t require.TestingT) {
		// given
		var expected *string
		var actual *string

		// when
		testutils.RequireBothNilOrEqual(t, expected, actual)
	}))

	t.Run("fails when expected pointer is nil but actual has a value", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		var expected *string
		actual := ptr.Addr("test")

		// when
		testutils.RequireBothNilOrEqual(t, expected, actual)
	}))

	t.Run("fails when expected has a value but actual pointer is nil", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		var actual *string

		// when
		testutils.RequireBothNilOrEqual(t, expected, actual)
	}))

	t.Run("fails when both expected and actual have values but they are different", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		actual := ptr.Addr("different")

		// when
		testutils.RequireBothNilOrEqual(t, expected, actual)
	}))
}

func Test_RequireBothNilOrEqualFn(t *testing.T) {
	haveTheSameLength := func(t require.TestingT, expected, actual string) {
		require.Len(t, actual, len(expected))
	}
	t.Run("passes when both pointers are non-nil and equal", testutils.RequireNoTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		actual := ptr.Addr("test")

		// when
		testutils.RequireBothNilOrEqualFn(t, expected, actual, haveTheSameLength)
	}))

	t.Run("passes when both pointers are nil", testutils.RequireNoTestFailure(func(t require.TestingT) {
		// given
		var expected *string
		var actual *string

		// when
		testutils.RequireBothNilOrEqualFn(t, expected, actual, haveTheSameLength)
	}))

	t.Run("fails when expected pointer is nil but actual has a value", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		var expected *string
		actual := ptr.Addr("test")

		// when
		testutils.RequireBothNilOrEqualFn(t, expected, actual, haveTheSameLength)
	}))

	t.Run("fails when expected has a value but actual pointer is nil", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		var actual *string

		// when
		testutils.RequireBothNilOrEqualFn(t, expected, actual, haveTheSameLength)
	}))

	t.Run("fails when both expected and actual have values but they are different", testutils.RequireTestFailure(func(t require.TestingT) {
		// given
		expected := ptr.Addr("test")
		actual := ptr.Addr("different")

		// when
		testutils.RequireBothNilOrEqualFn(t, expected, actual, haveTheSameLength)
	}))
}
