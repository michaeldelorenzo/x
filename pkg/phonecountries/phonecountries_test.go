package phonecountries_test

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/michaeldelorenzo/x/v3/pkg/phonecountries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allowedCountries = []string{
	"United States",
}

var expectedCountries = func() []phonecountries.Country {
	var want = []phonecountries.Country{
		{
			Name:         "Canada",
			OfficialName: "Canada",
			Alpha2Code:   "CA",
			CallingCodes: []string{"1"},
		},
		{
			Name:         "Puerto Rico",
			OfficialName: "Commonwealth of Puerto Rico",
			Alpha2Code:   "PR",
			CallingCodes: []string{"1787", "1939"},
		},
		{
			Name:         "United States",
			OfficialName: "United States of America",
			Alpha2Code:   "US",
			CallingCodes: []string{"1"},
		},
	}

	var countries = make(map[string]phonecountries.Country)
	var countryNames []string
	for _, c := range want {
		countryNames = append(countryNames, c.Name)

		sort.Strings(c.CallingCodes)
		countries[c.Name] = c
	}

	sort.Strings(countryNames)
	var expected []phonecountries.Country
	for _, name := range countryNames {
		expected = append(expected, countries[name])
	}

	return expected
}()

func TestAllowedCountries(t *testing.T) {
	out, err := phonecountries.AllowedCountries()
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
	assert.Equal(t, len(expectedCountries), len(out))
	assert.ElementsMatch(t, expectedCountries, out)
}

func TestAllowedCountriesMatch(t *testing.T) {
	ticketCountries := []struct {
		Name        string
		CallingCode string
	}{
		{Name: "Canada", CallingCode: "1"},
		{Name: "Puerto Rico", CallingCode: "1787"},
		{Name: "United States", CallingCode: "1"},
	}

	for _, expected := range ticketCountries {
		out, err := phonecountries.FindCountryByName(expected.Name)
		require.NoError(t, err)

		assert.Equal(t, expected.Name, out.Name)
		assert.Contains(t, out.CallingCodes, expected.CallingCode, fmt.Sprintf("calling codes are not included for %s", out))
	}
}

func TestFindCountryByName(t *testing.T) {
	for _, in := range allowedCountries {
		out, err := phonecountries.FindCountryByName(in)
		require.NoError(t, err)
		assert.Equal(t, in, out.Name)
		assert.NotEmpty(t, out.CallingCodes)
	}

	_, err := phonecountries.FindCountryByName("Afghanistan")
	require.Error(t, err)

	expectedError := &phonecountries.CountryNotAllowedError{}
	assert.True(t, errors.As(err, &expectedError))
}

func TestFindCountryByAlpha(t *testing.T) {
	tests := []struct {
		inputAlpha          string
		expectedCountryName string
	}{
		{
			inputAlpha:          "US",
			expectedCountryName: "United States",
		},
		{
			inputAlpha:          "CA",
			expectedCountryName: "Canada",
		},
		{
			inputAlpha:          "PR",
			expectedCountryName: "Puerto Rico",
		},
	}

	for _, tc := range tests {
		out, err := phonecountries.FindCountryByAlpha(tc.inputAlpha)
		require.NoError(t, err)

		assert.Equal(t, tc.expectedCountryName, out.Name)
		assert.Equal(t, tc.inputAlpha, out.Alpha2Code)
	}

	// country not allowed
	_, err := phonecountries.FindCountryByAlpha("AF")
	notAllowedErr := phonecountries.NewCountryNotAllowedError(
		phonecountries.Country{
			Name:       "Afghanistan",
			Alpha2Code: "AF",
		},
	)
	require.Error(t, err)
	require.ErrorAs(t, err, &notAllowedErr)

	// invalid country
	_, err = phonecountries.FindCountryByAlpha("XX")
	notFoundErr := phonecountries.NewCountryNotFoundError("XX")
	require.Error(t, err)
	require.ErrorAs(t, err, &notFoundErr)
}

func TestFindCountryByCallingCode(t *testing.T) {
	ticketCountries := []struct {
		Name        string
		CallingCode string
	}{
		{Name: "Canada", CallingCode: "1"},
		{Name: "Puerto Rico", CallingCode: "1787"},
		{Name: "United States", CallingCode: "1"},
	}

	for _, expected := range ticketCountries {
		code, _ := strconv.Atoi(expected.CallingCode)
		err := phonecountries.ValidateCountryCallingCode(int32(code))
		require.NoError(t, err)
	}

	err := phonecountries.ValidateCountryCallingCode(int32(93))
	require.Error(t, err)

	expectedError := &phonecountries.CountryCodeNotSupportedError{}
	assert.True(t, errors.As(err, &expectedError))
}
