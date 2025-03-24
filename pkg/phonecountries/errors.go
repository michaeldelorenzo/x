package phonecountries

import "fmt"

type CountryNotAllowedError struct {
	country Country
}

func (c *CountryNotAllowedError) Error() string {
	return fmt.Sprintf("`%s` is not an allowed country", c.country)
}

func NewCountryNotAllowedError(in Country) *CountryNotAllowedError {
	return &CountryNotAllowedError{country: in}
}

type CountryNotFoundError struct {
	search string
}

func (c *CountryNotFoundError) Error() string {
	return fmt.Sprintf("no country found using search `%s`", c.search)
}

func NewCountryNotFoundError(in string) *CountryNotFoundError {
	return &CountryNotFoundError{search: in}
}

type CountryCodeNotSupportedError struct {
	code int32
}

func (c *CountryCodeNotSupportedError) Error() string {
	return fmt.Sprintf("no allowable country found using code `%d`", c.code)
}

func NewCountryCodeNotSupportedError(in int32) *CountryCodeNotSupportedError {
	return &CountryCodeNotSupportedError{code: in}
}
