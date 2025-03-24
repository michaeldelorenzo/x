package phonecountries_test

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/koneksahealth/x/pkg/phonecountries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allowedCountries = []string{
	"American Samoa",
	"Anguilla",
	"Antigua and Barbuda",
	"Argentina",
	"Aruba",
	"Ascension",
	"Austria",
	"Bahamas",
	"Barbados",
	"Belarus",
	"Belize",
	"Bermuda",
	"Canada",
	"Caribbean Netherlands",
	"Cayman Islands",
	"Chile",
	"Costa Rica",
	"Croatia",
	"Cuba",
	"Curaçao",
	"Czech Republic",
	"Dominica",
	"Dominican Republic",
	"Dominican Republic",
	"Dominican Republic",
	"Dominican Republic",
	"El Salvador",
	"Finland",
	"France",
	"Germany",
	"Greece",
	"Greenland",
	"Grenada",
	"Guadeloupe",
	"Guam",
	"Guatemala",
	"Haiti",
	"Honduras",
	"Hungary",
	"India",
	"Israel",
	"Italy",
	"Jamaica",
	"Japan",
	"Kazakhstan",
	"Latvia",
	"Lithuania",
	"Martinique",
	"Mexico",
	"Montserrat",
	"Nicaragua",
	"Northern Mariana Islands",
	"Norway",
	"Panama",
	"Poland",
	"Portugal",
	"Puerto Rico",
	"Romania",
	"Russia",
	"Saint Kitts and Nevis",
	"Saint Lucia",
	"Saint Pierre and Miquelon",
	"Saint Vincent Grenadines",
	"Spain",
	"Sweden",
	"Switzerland",
	"Trinidad and Tobago",
	"Turkey",
	"Turks and Caicos Islands",
	"Ukraine",
	"United Kingdom",
	"United States",
	"Virgin Islands, British",
	"Virgin Islands, U.S.",
}

var expectedCountries = func() []phonecountries.Country {
	var want = []phonecountries.Country{
		{
			Name:         "American Samoa",
			OfficialName: "American Samoa",
			Alpha2Code:   "AS",
			CallingCodes: []string{"1684"},
		},
		{
			Name:         "Anguilla",
			OfficialName: "Anguilla",
			Alpha2Code:   "AI",
			CallingCodes: []string{"1264"},
		},
		{
			Name:         "Antigua and Barbuda",
			OfficialName: "Antigua and Barbuda",
			Alpha2Code:   "AG",
			CallingCodes: []string{"1268"},
		},
		{
			Name:         "Argentina",
			OfficialName: "Argentine Republic",
			Alpha2Code:   "AR",
			CallingCodes: []string{"54"},
		},
		{
			Name:         "Aruba",
			OfficialName: "Aruba",
			Alpha2Code:   "AW",
			CallingCodes: []string{"297"},
		},
		{
			Name:         "Ascension",
			OfficialName: "Saint Helena, Ascension and Tristan da Cunha",
			Alpha2Code:   "SH",
			CallingCodes: []string{"290", "247"},
		},
		{
			Name:         "Austria",
			OfficialName: "Republic of Austria",
			Alpha2Code:   "AT",
			CallingCodes: []string{"43"},
		},
		{
			Name:         "Bahamas",
			OfficialName: "Commonwealth of the Bahamas",
			Alpha2Code:   "BS",
			CallingCodes: []string{"1242"},
		},
		{
			Name:         "Barbados",
			OfficialName: "Barbados",
			Alpha2Code:   "BB",
			CallingCodes: []string{"1246"},
		},
		{
			Name:         "Belarus",
			OfficialName: "Republic of Belarus",
			Alpha2Code:   "BY",
			CallingCodes: []string{"375"},
		},
		{
			Name:         "Belize",
			OfficialName: "Belize",
			Alpha2Code:   "BZ",
			CallingCodes: []string{"501"},
		},
		{
			Name:         "Bermuda",
			OfficialName: "Bermuda",
			Alpha2Code:   "BM",
			CallingCodes: []string{"1441"},
		},
		{
			Name:         "Canada",
			OfficialName: "Canada",
			Alpha2Code:   "CA",
			CallingCodes: []string{"1"},
		},
		{
			Name:         "Caribbean Netherlands",
			OfficialName: "Bonaire, Sint Eustatius and Saba",
			Alpha2Code:   "BQ",
			CallingCodes: []string{"599"},
		},
		{
			Name:         "Cayman Islands",
			OfficialName: "Cayman Islands",
			Alpha2Code:   "KY",
			CallingCodes: []string{"1345"},
		},
		{
			Name:         "Chile",
			OfficialName: "Republic of Chile",
			Alpha2Code:   "CL",
			CallingCodes: []string{"56"},
		},
		{
			Name:         "Costa Rica",
			OfficialName: "Republic of Costa Rica",
			Alpha2Code:   "CR",
			CallingCodes: []string{"506"},
		},
		{
			Name:         "Croatia",
			OfficialName: "Republic of Croatia",
			Alpha2Code:   "HR",
			CallingCodes: []string{"385"},
		},
		{
			Name:         "Cuba",
			OfficialName: "Republic of Cuba",
			Alpha2Code:   "CU",
			CallingCodes: []string{"53"},
		},
		{
			Name:         "Curaçao",
			OfficialName: "Country of Curaçao",
			Alpha2Code:   "CW",
			CallingCodes: []string{"5999", "599"},
		},
		{
			Name:         "Czech Republic",
			OfficialName: "Czech Republic",
			Alpha2Code:   "CZ",
			CallingCodes: []string{"420"},
		},
		{
			Name:         "Dominica",
			OfficialName: "Commonwealth of Dominica",
			Alpha2Code:   "DM",
			CallingCodes: []string{"1767"},
		},
		{
			Name:         "Dominican Republic",
			OfficialName: "Dominican Republic",
			Alpha2Code:   "DO",
			CallingCodes: []string{"1809", "1829", "1849", "1809201"},
		},
		{
			Name:         "El Salvador",
			OfficialName: "Republic of El Salvador",
			Alpha2Code:   "SV",
			CallingCodes: []string{"503"},
		},
		{
			Name:         "Finland",
			OfficialName: "Republic of Finland",
			Alpha2Code:   "FI",
			CallingCodes: []string{"358"},
		},
		{
			Name:         "France",
			OfficialName: "French Republic",
			Alpha2Code:   "FR",
			CallingCodes: []string{"33"},
		},
		{
			Name:         "Germany",
			OfficialName: "Federal Republic of Germany",
			Alpha2Code:   "DE",
			CallingCodes: []string{"49"},
		},
		{
			Name:         "Greece",
			OfficialName: "Hellenic Republic",
			Alpha2Code:   "GR",
			CallingCodes: []string{"30"},
		},
		{
			Name:         "Greenland",
			OfficialName: "Greenland",
			Alpha2Code:   "GL",
			CallingCodes: []string{"299"},
		},
		{
			Name:         "Grenada",
			OfficialName: "Grenada",
			Alpha2Code:   "GD",
			CallingCodes: []string{"1473"},
		},
		{
			Name:         "Guadeloupe",
			OfficialName: "Guadeloupe",
			Alpha2Code:   "GP",
			CallingCodes: []string{"590"},
		},
		{
			Name:         "Guam",
			OfficialName: "Guam",
			Alpha2Code:   "GU",
			CallingCodes: []string{"1671"},
		},
		{
			Name:         "Guatemala",
			OfficialName: "Republic of Guatemala",
			Alpha2Code:   "GT",
			CallingCodes: []string{"502"},
		},
		{
			Name:         "Haiti",
			OfficialName: "Republic of Haiti",
			Alpha2Code:   "HT",
			CallingCodes: []string{"509"},
		},
		{
			Name:         "Honduras",
			OfficialName: "Republic of Honduras",
			Alpha2Code:   "HN",
			CallingCodes: []string{"504"},
		},
		{
			Name:         "Hungary",
			OfficialName: "Hungary",
			Alpha2Code:   "HU",
			CallingCodes: []string{"36"},
		},
		{
			Name:         "India",
			OfficialName: "Republic of India",
			Alpha2Code:   "IN",
			CallingCodes: []string{"91"},
		},
		{
			Name:         "Israel",
			OfficialName: "State of Israel",
			Alpha2Code:   "IL",
			CallingCodes: []string{"972"},
		},
		{
			Name:         "Italy",
			OfficialName: "Italian Republic",
			Alpha2Code:   "IT",
			CallingCodes: []string{"39"},
		},
		{
			Name:         "Jamaica",
			OfficialName: "Jamaica",
			Alpha2Code:   "JM",
			CallingCodes: []string{"1876"},
		},
		{
			Name:         "Japan",
			OfficialName: "Japan",
			Alpha2Code:   "JP",
			CallingCodes: []string{"81"},
		},
		{
			Name:         "Kazakhstan",
			OfficialName: "Republic of Kazakhstan",
			Alpha2Code:   "KZ",
			CallingCodes: []string{"76", "77", "7"},
		},
		{
			Name:         "Latvia",
			OfficialName: "Republic of Latvia",
			Alpha2Code:   "LV",
			CallingCodes: []string{"371"},
		},
		{
			Name:         "Lithuania",
			OfficialName: "Republic of Lithuania",
			Alpha2Code:   "LT",
			CallingCodes: []string{"370"},
		},
		{
			Name:         "Martinique",
			OfficialName: "Martinique",
			Alpha2Code:   "MQ",
			CallingCodes: []string{"596"},
		},
		{
			Name:         "Mexico",
			OfficialName: "United Mexican States",
			Alpha2Code:   "MX",
			CallingCodes: []string{"52"},
		},
		{
			Name:         "Montserrat",
			OfficialName: "Montserrat",
			Alpha2Code:   "MS",
			CallingCodes: []string{"1664"},
		},
		{
			Name:         "Nicaragua",
			OfficialName: "Republic of Nicaragua",
			Alpha2Code:   "NI",
			CallingCodes: []string{"505"},
		},
		{
			Name:         "Northern Mariana Islands",
			OfficialName: "Commonwealth of the Northern Mariana Islands",
			Alpha2Code:   "MP",
			CallingCodes: []string{"1670"},
		},
		{
			Name:         "Norway",
			OfficialName: "Kingdom of Norway",
			Alpha2Code:   "NO",
			CallingCodes: []string{"47"},
		},
		{
			Name:         "Panama",
			OfficialName: "Republic of Panama",
			Alpha2Code:   "PA",
			CallingCodes: []string{"507"},
		},
		{
			Name:         "Poland",
			OfficialName: "Republic of Poland",
			Alpha2Code:   "PL",
			CallingCodes: []string{"48"},
		},
		{
			Name:         "Portugal",
			OfficialName: "Portuguese Republic",
			Alpha2Code:   "PT",
			CallingCodes: []string{"351"},
		},
		{
			Name:         "Puerto Rico",
			OfficialName: "Commonwealth of Puerto Rico",
			Alpha2Code:   "PR",
			CallingCodes: []string{"1787", "1939"},
		},
		{
			Name:         "Romania",
			OfficialName: "Romania",
			Alpha2Code:   "RO",
			CallingCodes: []string{"40"},
		},
		{
			Name:         "Russia",
			OfficialName: "Russian Federation",
			Alpha2Code:   "RU",
			CallingCodes: []string{"7"},
		},
		{
			Name:         "Saint Kitts and Nevis",
			OfficialName: "Federation of Saint Christopher and Nevisa",
			Alpha2Code:   "KN",
			CallingCodes: []string{"1869"},
		},
		{
			Name:         "Saint Lucia",
			OfficialName: "Saint Lucia",
			Alpha2Code:   "LC",
			CallingCodes: []string{"1758"},
		},
		{
			Name:         "Saint Pierre and Miquelon",
			OfficialName: "Saint Pierre and Miquelon",
			Alpha2Code:   "PM",
			CallingCodes: []string{"508"},
		},
		{
			Name:         "Saint Vincent Grenadines",
			OfficialName: "Saint Vincent and the Grenadines",
			Alpha2Code:   "VC",
			CallingCodes: []string{"1784"},
		},
		{
			Name:         "Spain",
			OfficialName: "Kingdom of Spain",
			Alpha2Code:   "ES",
			CallingCodes: []string{"34"},
		},
		{
			Name:         "Sweden",
			OfficialName: "Kingdom of Sweden",
			Alpha2Code:   "SE",
			CallingCodes: []string{"46"},
		},
		{
			Name:         "Switzerland",
			OfficialName: "Swiss Confederation",
			Alpha2Code:   "CH",
			CallingCodes: []string{"41"},
		},
		{
			Name:         "Trinidad and Tobago",
			OfficialName: "Republic of Trinidad and Tobago",
			Alpha2Code:   "TT",
			CallingCodes: []string{"1868"},
		},
		{
			Name:         "Turkey",
			OfficialName: "Republic of Turkey",
			Alpha2Code:   "TR",
			CallingCodes: []string{"90"},
		},
		{
			Name:         "Turks and Caicos Islands",
			OfficialName: "Turks and Caicos Islands",
			Alpha2Code:   "TC",
			CallingCodes: []string{"1649"},
		},
		{
			Name:         "Ukraine",
			OfficialName: "Ukraine",
			Alpha2Code:   "UA",
			CallingCodes: []string{"380"},
		},
		{
			Name:         "United Kingdom",
			OfficialName: "United Kingdom of Great Britain and Northern Ireland",
			Alpha2Code:   "GB",
			CallingCodes: []string{"44"},
		},
		{
			Name:         "United States",
			OfficialName: "United States of America",
			Alpha2Code:   "US",
			CallingCodes: []string{"1"},
		},
		{
			Name:         "Virgin Islands, British",
			OfficialName: "Virgin Islands",
			Alpha2Code:   "VG",
			CallingCodes: []string{"1284"},
		},
		{
			Name:         "Virgin Islands, U.S.",
			OfficialName: "Virgin Islands of the United States",
			Alpha2Code:   "VI",
			CallingCodes: []string{"1340"},
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

func TestAllowedCountriesMatchKH_5982(t *testing.T) {
	// Countries and codes were listed on https://koneksa-health.atlassian.net/browse/KH-5982
	ticketCountries := []struct {
		Name        string
		CallingCode string
	}{
		{Name: "American Samoa", CallingCode: "1684"},
		{Name: "Anguilla", CallingCode: "1264"},
		{Name: "Antigua and Barbuda", CallingCode: "1268"},
		{Name: "Argentina", CallingCode: "54"},
		{Name: "Aruba", CallingCode: "297"},
		{Name: "Ascension", CallingCode: "247"},
		{Name: "Austria", CallingCode: "43"},
		{Name: "Bahamas", CallingCode: "1242"},
		{Name: "Barbados", CallingCode: "1246"},
		{Name: "Belarus", CallingCode: "375"},
		{Name: "Belize", CallingCode: "501"},
		{Name: "Bermuda", CallingCode: "1441"},
		{Name: "Canada", CallingCode: "1"},
		{Name: "Caribbean Netherlands", CallingCode: "599"},
		{Name: "Cayman Islands", CallingCode: "1345"},
		{Name: "Chile", CallingCode: "56"},
		{Name: "Costa Rica", CallingCode: "506"},
		{Name: "Croatia", CallingCode: "385"},
		{Name: "Cuba", CallingCode: "53"},
		{Name: "Curaçao", CallingCode: "599"},
		{Name: "Czech Republic", CallingCode: "420"},
		{Name: "Dominica", CallingCode: "1767"},
		{Name: "Dominican Republic", CallingCode: "1809"},
		{Name: "Dominican Republic", CallingCode: "1809201"},
		{Name: "Dominican Republic", CallingCode: "1829"},
		{Name: "Dominican Republic", CallingCode: "1849"},
		{Name: "El Salvador", CallingCode: "503"},
		{Name: "Finland", CallingCode: "358"},
		{Name: "France", CallingCode: "33"},
		{Name: "Germany", CallingCode: "49"},
		{Name: "Greece", CallingCode: "30"},
		{Name: "Greenland", CallingCode: "299"},
		{Name: "Grenada", CallingCode: "1473"},
		{Name: "Guadeloupe", CallingCode: "590"},
		{Name: "Guam", CallingCode: "1671"},
		{Name: "Guatemala", CallingCode: "502"},
		{Name: "Haiti", CallingCode: "509"},
		{Name: "Honduras", CallingCode: "504"},
		{Name: "Hungary", CallingCode: "36"},
		{Name: "India", CallingCode: "91"},
		{Name: "Israel", CallingCode: "972"},
		{Name: "Italy", CallingCode: "39"},
		{Name: "Jamaica", CallingCode: "1876"},
		{Name: "Japan", CallingCode: "81"},
		{Name: "Kazakhstan", CallingCode: "7"},
		{Name: "Latvia", CallingCode: "371"},
		{Name: "Lithuania", CallingCode: "370"},
		{Name: "Martinique", CallingCode: "596"},
		{Name: "Mexico", CallingCode: "52"},
		{Name: "Montserrat", CallingCode: "1664"},
		{Name: "Nicaragua", CallingCode: "505"},
		{Name: "Northern Mariana Islands", CallingCode: "1670"},
		{Name: "Norway", CallingCode: "47"},
		{Name: "Panama", CallingCode: "507"},
		{Name: "Poland", CallingCode: "48"},
		{Name: "Portugal", CallingCode: "351"},
		{Name: "Puerto Rico", CallingCode: "1787"},
		{Name: "Romania", CallingCode: "40"},
		{Name: "Russia", CallingCode: "7"},
		{Name: "Saint Kitts and Nevis", CallingCode: "1869"},
		{Name: "Saint Lucia", CallingCode: "1758"},
		{Name: "Saint Pierre and Miquelon", CallingCode: "508"},
		{Name: "Saint Vincent Grenadines", CallingCode: "1784"},
		{Name: "Spain", CallingCode: "34"},
		{Name: "Sweden", CallingCode: "46"},
		{Name: "Switzerland", CallingCode: "41"},
		{Name: "Trinidad and Tobago", CallingCode: "1868"},
		{Name: "Turkey", CallingCode: "90"},
		{Name: "Turks and Caicos Islands", CallingCode: "1649"},
		{Name: "Ukraine", CallingCode: "380"},
		{Name: "United Kingdom", CallingCode: "44"},
		{Name: "United States", CallingCode: "1"},
		{Name: "Virgin Islands, British", CallingCode: "1284"},
		{Name: "Virgin Islands, U.S.", CallingCode: "1340"},
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
			inputAlpha:          "UA",
			expectedCountryName: "Ukraine",
		},
		{
			inputAlpha:          "CZ",
			expectedCountryName: "Czech Republic",
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
		{Name: "American Samoa", CallingCode: "1684"},
		{Name: "Anguilla", CallingCode: "1264"},
		{Name: "Antigua and Barbuda", CallingCode: "1268"},
		{Name: "Argentina", CallingCode: "54"},
		{Name: "Aruba", CallingCode: "297"},
		{Name: "Ascension", CallingCode: "247"},
		{Name: "Austria", CallingCode: "43"},
		{Name: "Bahamas", CallingCode: "1242"},
		{Name: "Barbados", CallingCode: "1246"},
		{Name: "Belarus", CallingCode: "375"},
		{Name: "Belize", CallingCode: "501"},
		{Name: "Bermuda", CallingCode: "1441"},
		{Name: "Canada", CallingCode: "1"},
		{Name: "Caribbean Netherlands", CallingCode: "599"},
		{Name: "Cayman Islands", CallingCode: "1345"},
		{Name: "Chile", CallingCode: "56"},
		{Name: "Costa Rica", CallingCode: "506"},
		{Name: "Croatia", CallingCode: "385"},
		{Name: "Cuba", CallingCode: "53"},
		{Name: "Curaçao", CallingCode: "599"},
		{Name: "Czech Republic", CallingCode: "420"},
		{Name: "Dominica", CallingCode: "1767"},
		{Name: "Dominican Republic", CallingCode: "1809"},
		{Name: "Dominican Republic", CallingCode: "1809201"},
		{Name: "Dominican Republic", CallingCode: "1829"},
		{Name: "Dominican Republic", CallingCode: "1849"},
		{Name: "El Salvador", CallingCode: "503"},
		{Name: "Finland", CallingCode: "358"},
		{Name: "France", CallingCode: "33"},
		{Name: "Germany", CallingCode: "49"},
		{Name: "Greece", CallingCode: "30"},
		{Name: "Greenland", CallingCode: "299"},
		{Name: "Grenada", CallingCode: "1473"},
		{Name: "Guadeloupe", CallingCode: "590"},
		{Name: "Guam", CallingCode: "1671"},
		{Name: "Guatemala", CallingCode: "502"},
		{Name: "Haiti", CallingCode: "509"},
		{Name: "Honduras", CallingCode: "504"},
		{Name: "Hungary", CallingCode: "36"},
		{Name: "India", CallingCode: "91"},
		{Name: "Israel", CallingCode: "972"},
		{Name: "Italy", CallingCode: "39"},
		{Name: "Jamaica", CallingCode: "1876"},
		{Name: "Japan", CallingCode: "81"},
		{Name: "Kazakhstan", CallingCode: "7"},
		{Name: "Latvia", CallingCode: "371"},
		{Name: "Lithuania", CallingCode: "370"},
		{Name: "Martinique", CallingCode: "596"},
		{Name: "Mexico", CallingCode: "52"},
		{Name: "Montserrat", CallingCode: "1664"},
		{Name: "Nicaragua", CallingCode: "505"},
		{Name: "Northern Mariana Islands", CallingCode: "1670"},
		{Name: "Norway", CallingCode: "47"},
		{Name: "Panama", CallingCode: "507"},
		{Name: "Poland", CallingCode: "48"},
		{Name: "Portugal", CallingCode: "351"},
		{Name: "Puerto Rico", CallingCode: "1787"},
		{Name: "Romania", CallingCode: "40"},
		{Name: "Russia", CallingCode: "7"},
		{Name: "Saint Kitts and Nevis", CallingCode: "1869"},
		{Name: "Saint Lucia", CallingCode: "1758"},
		{Name: "Saint Pierre and Miquelon", CallingCode: "508"},
		{Name: "Saint Vincent Grenadines", CallingCode: "1784"},
		{Name: "Spain", CallingCode: "34"},
		{Name: "Sweden", CallingCode: "46"},
		{Name: "Switzerland", CallingCode: "41"},
		{Name: "Trinidad and Tobago", CallingCode: "1868"},
		{Name: "Turkey", CallingCode: "90"},
		{Name: "Turks and Caicos Islands", CallingCode: "1649"},
		{Name: "Ukraine", CallingCode: "380"},
		{Name: "United Kingdom", CallingCode: "44"},
		{Name: "United States", CallingCode: "1"},
		{Name: "Virgin Islands, British", CallingCode: "1284"},
		{Name: "Virgin Islands, U.S.", CallingCode: "1340"},
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
