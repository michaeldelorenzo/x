// Package phonecountries contains the platform's source-of-truth list of supported countries for use with phone numbers.
package phonecountries

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/michaeldelorenzo/x/v3/pkg/utils/sequence"
	"github.com/pariz/gountries"
)

var countriesQuery = gountries.New()

// countryDisplayNameOverrides overrides the default country names to enable proper lookup and validation.
var countryDisplayNameOverrides = map[string]string{
	// mapping is:
	// common name: Preferred Name
	//"Saint Helena":                     "Ascension",
	//"Saint Vincent and the Grenadines": "Saint Vincent Grenadines",
	//"British Virgin Islands":           "Virgin Islands, British",
	//"United States Virgin Islands":     "Virgin Islands, U.S.",
}

// countryQueryNameOverrides inverts the countryDisplayNameOverrides to support querying.
var countryQueryNameOverrides = func() map[string]string {
	m := map[string]string{}
	for k, v := range countryDisplayNameOverrides {
		m[v] = k
	}
	return m
}()

// additionalCountryCallingCodes some countries have additional calling country codes that are not included in the gountries package.
var additionalCountryCallingCodes = map[string][]string{
	//"Curaçao":            {"599"},
	//"Dominican Republic": {"1809201"},
	//"Kazakhstan":         {"7"},
}

var allowedCountries = []string{
	"Canada",
	"Puerto Rico",
	"United States",
}

// Country represents a country and its properties.
type Country struct {
	// Name is the country's common name.
	Name         string
	OfficialName string
	Alpha2Code   string
	CallingCodes []string
}

func (c Country) String() string {
	return fmt.Sprintf("%s [%s]", c.Name, c.Alpha2Code)
}

// toCountryOutput translates the gountries.Country input to a phonecountries.Country output.
func toCountryOutput(in gountries.Country) Country {
	commonName := in.Name.Common
	if nameOverride, ok := countryDisplayNameOverrides[in.Name.Common]; ok {
		commonName = nameOverride
	}

	callingCodes := in.CallingCodes
	if additionalCodes, ok := additionalCountryCallingCodes[commonName]; ok {
		callingCodes = append(callingCodes, additionalCodes...)
	}
	// sorting for some predictability
	sort.Strings(callingCodes)

	return Country{
		Name:         commonName,
		OfficialName: in.Name.Official,
		Alpha2Code:   in.Alpha2,
		CallingCodes: callingCodes,
	}
}

// validate checks that the country is included in the allowed list.
func validate(input Country) (Country, error) {
	allowed, _ := AllowedCountries()

	for _, country := range allowed {
		if reflect.DeepEqual(country, input) {
			return input, nil
		}
	}

	return Country{}, NewCountryNotAllowedError(input)
}

// toCountryNameQuery accepts the query input and applies the country name override, if necessary.
func toCountryNameQuery(q string) string {
	if queryOverride, ok := countryQueryNameOverrides[q]; ok {
		return queryOverride
	}
	return q
}

// AllowedCountries returns the list of allowable countries.
func AllowedCountries() ([]Country, error) {
	var countries = make(map[string]Country)
	var countryNames []string
	for _, allowedCountryName := range allowedCountries {
		c, err := countriesQuery.FindCountryByName(toCountryNameQuery(allowedCountryName))
		if err != nil {
			return nil, NewCountryNotFoundError(allowedCountryName)
		}

		country := toCountryOutput(c)
		countries[country.Name] = country

		countryNames = append(countryNames, country.Name)
	}

	// sort the country names
	sort.Strings(countryNames)

	var out []Country
	for _, name := range countryNames {
		out = append(out, countries[name])
	}
	return out, nil
}

// FindCountryByName searches the allowable country list for the country by name.
func FindCountryByName(in string) (Country, error) {
	c, err := countriesQuery.FindCountryByName(toCountryNameQuery(in))
	if err != nil {
		return Country{}, NewCountryNotFoundError(in)
	}

	return validate(toCountryOutput(c))
}

// FindCountryByAlpha searches the allowable country list for the country using the provided alpha-code.
func FindCountryByAlpha(in string) (Country, error) {
	c, err := countriesQuery.FindCountryByAlpha(in)
	if err != nil {
		return Country{}, NewCountryNotFoundError(in)
	}

	return validate(toCountryOutput(c))
}

// ValidateCountryCallingCode searches the allowable country list for the country using the provided calling code. Returns
// no error if the calling code is permitted.
func ValidateCountryCallingCode(in int32) error {
	allowed, _ := AllowedCountries()
	for _, c := range allowed {
		if sequence.Contains(c.CallingCodes, fmt.Sprintf("%d", in)) {
			return nil
		}
	}

	return NewCountryCodeNotSupportedError(in)
}
