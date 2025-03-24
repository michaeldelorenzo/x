package ptr_test

import (
	"testing"

	"github.com/koneksahealth/x/pkg/utils/ptr"
	"github.com/stretchr/testify/require"
)

func TestPtr_Addr(t *testing.T) {
	t.Run("supports taking the address of a primative", func(t *testing.T) {
		// given
		input := "some value"

		// when
		output := ptr.Addr("some value")

		// then
		require.NotNil(t, output)
		require.EqualValues(t, input, *output)
	})
}

func TestPtr_Deref(t *testing.T) {
	t.Run("dereferences the pointer when it is not nil", func(t *testing.T) {
		// given
		expected := "some value"
		input := &expected

		// when
		output := ptr.Deref(input)

		// then
		require.NotEmpty(t, output)
		require.EqualValues(t, expected, output)
	})

	t.Run("returns the zero value for the type when the pointer is nil", func(t *testing.T) {
		// given
		var input *string

		// when
		output := ptr.Deref(input)

		// then
		require.Empty(t, output)
	})
}

func TestPtr_DerefWithFallback(t *testing.T) {
	t.Run("dereferences the pointer when it is not nil", func(t *testing.T) {
		// given
		expected := "some value"
		input := &expected

		// when
		output := ptr.DerefWithFallback(input, "another fallback")

		// then
		require.NotEmpty(t, output)
		require.EqualValues(t, expected, output)
	})

	t.Run("returns the fallback value when the pointer is nil", func(t *testing.T) {
		// given
		expected := "my fallback value"
		var input *string

		// when
		output := ptr.DerefWithFallback(input, expected)

		// then
		require.NotEmpty(t, output)
		require.EqualValues(t, expected, output)
	})
}
