package locales_test

import (
	"testing"

	"github.com/koneksahealth/x/pkg/locales"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestLocales_Validate(t *testing.T) {
	t.Run("returns true if the input is a supported BCP 47 locale code", func(t *testing.T) {
		output := locales.Validate("en-US")
		require.True(t, output)
	})

	t.Run("returns false if the input is not a supported BCP 47 locale code", func(t *testing.T) {
		output := locales.Validate("nb-no") // non-canonical language, should be no-no
		require.False(t, output)
	})

	t.Run("returns false if the input not a valid BCP 47 locale code", func(t *testing.T) {
		output := locales.Validate("bad-stuff")
		require.False(t, output)
	})
}

var expectedLocales = []locales.Locale{
	{
		Code:                "en-US",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "United States",
		RegionBCP47Code:     "US",
	},
	{
		Code:                "en-GB",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "United Kingdom",
		RegionBCP47Code:     "GB",
	},
	{
		Code:                "af-ZA",
		LanguageDisplayName: "Afrikaans",
		LanguageBCP47Code:   "af",
		RegionDisplayName:   "South Africa",
		RegionBCP47Code:     "ZA",
	},
	{
		Code:                "ar-EG",
		LanguageDisplayName: "Arabic",
		LanguageBCP47Code:   "ar",
		RegionDisplayName:   "Egypt",
		RegionBCP47Code:     "EG",
	},
	{
		Code:                "ar-IL",
		LanguageDisplayName: "Arabic",
		LanguageBCP47Code:   "ar",
		RegionDisplayName:   "Israel",
		RegionBCP47Code:     "IL",
	},
	{
		Code:                "ar-LB",
		LanguageDisplayName: "Arabic",
		LanguageBCP47Code:   "ar",
		RegionDisplayName:   "Lebanon",
		RegionBCP47Code:     "LB",
	},
	{
		Code:                "ar-QA",
		LanguageDisplayName: "Arabic",
		LanguageBCP47Code:   "ar",
		RegionDisplayName:   "Qatar",
		RegionBCP47Code:     "QA",
	},
	{
		Code:                "ar-AE",
		LanguageDisplayName: "Arabic",
		LanguageBCP47Code:   "ar",
		RegionDisplayName:   "United Arab Emirates",
		RegionBCP47Code:     "AE",
	},
	{
		Code:                "be-BY",
		LanguageDisplayName: "Belarusian",
		LanguageBCP47Code:   "be",
		RegionDisplayName:   "Belarus",
		RegionBCP47Code:     "BY",
	},
	{
		Code:                "bn-IN",
		LanguageDisplayName: "Bangla",
		LanguageBCP47Code:   "bn",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "bg-BG",
		LanguageDisplayName: "Bulgarian",
		LanguageBCP47Code:   "bg",
		RegionDisplayName:   "Bulgaria",
		RegionBCP47Code:     "BG",
	},
	{
		Code:                "zh-Hans-CN",
		LanguageDisplayName: "Chinese",
		LanguageBCP47Code:   "zh",
		RegionDisplayName:   "China",
		RegionBCP47Code:     "CN",
	},
	{
		Code:                "zh-Hans-MY",
		LanguageDisplayName: "Chinese",
		LanguageBCP47Code:   "zh",
		RegionDisplayName:   "Malaysia",
		RegionBCP47Code:     "MY",
	},
	{
		Code:                "zh-Hans-SG",
		LanguageDisplayName: "Chinese",
		LanguageBCP47Code:   "zh",
		RegionDisplayName:   "Singapore",
		RegionBCP47Code:     "SG",
	},
	{
		Code:                "hr-HR",
		LanguageDisplayName: "Croatian",
		LanguageBCP47Code:   "hr",
		RegionDisplayName:   "Croatia",
		RegionBCP47Code:     "HR",
	},
	{
		Code:                "cs-CZ",
		LanguageDisplayName: "Czech",
		LanguageBCP47Code:   "cs",
		RegionDisplayName:   "Czechia",
		RegionBCP47Code:     "CZ",
	},
	{
		Code:                "da-DK",
		LanguageDisplayName: "Danish",
		LanguageBCP47Code:   "da",
		RegionDisplayName:   "Denmark",
		RegionBCP47Code:     "DK",
	},
	{
		Code:                "nl-BE",
		LanguageDisplayName: "Dutch",
		LanguageBCP47Code:   "nl",
		RegionDisplayName:   "Belgium",
		RegionBCP47Code:     "BE",
	},
	{
		Code:                "nl-NL",
		LanguageDisplayName: "Dutch",
		LanguageBCP47Code:   "nl",
		RegionDisplayName:   "Netherlands",
		RegionBCP47Code:     "NL",
	},
	{
		Code:                "en-AU",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Australia",
		RegionBCP47Code:     "AU",
	},
	{
		Code:                "en-CA",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Canada",
		RegionBCP47Code:     "CA",
	},
	{
		Code:                "en-IN",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "en-IE",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Ireland",
		RegionBCP47Code:     "IE",
	},
	{
		Code:                "en-KE",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Kenya",
		RegionBCP47Code:     "KE",
	},
	{
		Code:                "en-MY",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Malaysia",
		RegionBCP47Code:     "MY",
	},
	{
		Code:                "en-SG",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "Singapore",
		RegionBCP47Code:     "SG",
	},
	{
		Code:                "en-ZA",
		LanguageDisplayName: "English",
		LanguageBCP47Code:   "en",
		RegionDisplayName:   "South Africa",
		RegionBCP47Code:     "ZA",
	},
	{
		Code:                "et-EE",
		LanguageDisplayName: "Estonian",
		LanguageBCP47Code:   "et",
		RegionDisplayName:   "Estonia",
		RegionBCP47Code:     "EE",
	},
	{
		Code:                "fi-FI",
		LanguageDisplayName: "Finnish",
		LanguageBCP47Code:   "fi",
		RegionDisplayName:   "Finland",
		RegionBCP47Code:     "FI",
	},
	{
		Code:                "sv-FI",
		LanguageDisplayName: "Swedish",
		LanguageBCP47Code:   "sv",
		RegionDisplayName:   "Finland",
		RegionBCP47Code:     "FI",
	},
	{
		Code:                "fr-BE",
		LanguageDisplayName: "French",
		LanguageBCP47Code:   "fr",
		RegionDisplayName:   "Belgium",
		RegionBCP47Code:     "BE",
	},
	{
		Code:                "fr-CA",
		LanguageDisplayName: "French",
		LanguageBCP47Code:   "fr",
		RegionDisplayName:   "Canada",
		RegionBCP47Code:     "CA",
	},
	{
		Code:                "fr-FR",
		LanguageDisplayName: "French",
		LanguageBCP47Code:   "fr",
		RegionDisplayName:   "France",
		RegionBCP47Code:     "FR",
	},
	{
		Code:                "fr-CH",
		LanguageDisplayName: "French",
		LanguageBCP47Code:   "fr",
		RegionDisplayName:   "Switzerland",
		RegionBCP47Code:     "CH",
	},
	{
		Code:                "ka-GE",
		LanguageDisplayName: "Georgian",
		LanguageBCP47Code:   "ka",
		RegionDisplayName:   "Georgia",
		RegionBCP47Code:     "GE",
	},
	{
		Code:                "de-AT",
		LanguageDisplayName: "German",
		LanguageBCP47Code:   "de",
		RegionDisplayName:   "Austria",
		RegionBCP47Code:     "AT",
	},
	{
		Code:                "de-BE",
		LanguageDisplayName: "German",
		LanguageBCP47Code:   "de",
		RegionDisplayName:   "Belgium",
		RegionBCP47Code:     "BE",
	},
	{
		Code:                "de-DE",
		LanguageDisplayName: "German",
		LanguageBCP47Code:   "de",
		RegionDisplayName:   "Germany",
		RegionBCP47Code:     "DE",
	},
	{
		Code:                "de-CH",
		LanguageDisplayName: "German",
		LanguageBCP47Code:   "de",
		RegionDisplayName:   "Switzerland",
		RegionBCP47Code:     "CH",
	},
	{
		Code:                "el-GR",
		LanguageDisplayName: "Greek",
		LanguageBCP47Code:   "el",
		RegionDisplayName:   "Greece",
		RegionBCP47Code:     "GR",
	},
	{
		Code:                "gu-IN",
		LanguageDisplayName: "Gujarati",
		LanguageBCP47Code:   "gu",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "he-IL",
		LanguageDisplayName: "Hebrew",
		LanguageBCP47Code:   "he",
		RegionDisplayName:   "Israel",
		RegionBCP47Code:     "IL",
	},
	{
		Code:                "hi-IN",
		LanguageDisplayName: "Hindi",
		LanguageBCP47Code:   "hi",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "hu-HU",
		LanguageDisplayName: "Hungarian",
		LanguageBCP47Code:   "hu",
		RegionDisplayName:   "Hungary",
		RegionBCP47Code:     "HU",
	},
	{
		Code:                "ga-IE",
		LanguageDisplayName: "Irish",
		LanguageBCP47Code:   "ga",
		RegionDisplayName:   "Ireland",
		RegionBCP47Code:     "IE",
	},
	{
		Code:                "it-IT",
		LanguageDisplayName: "Italian",
		LanguageBCP47Code:   "it",
		RegionDisplayName:   "Italy",
		RegionBCP47Code:     "IT",
	},
	{
		Code:                "it-CH",
		LanguageDisplayName: "Italian",
		LanguageBCP47Code:   "it",
		RegionDisplayName:   "Switzerland",
		RegionBCP47Code:     "CH",
	},
	{
		Code:                "ja-JP",
		LanguageDisplayName: "Japanese",
		LanguageBCP47Code:   "ja",
		RegionDisplayName:   "Japan",
		RegionBCP47Code:     "JP",
	},
	{
		Code:                "kn-IN",
		LanguageDisplayName: "Kannada",
		LanguageBCP47Code:   "kn",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "ko-KR",
		LanguageDisplayName: "Korean",
		LanguageBCP47Code:   "ko",
		RegionDisplayName:   "South Korea",
		RegionBCP47Code:     "KR",
	},
	{
		Code:                "lv-LV",
		LanguageDisplayName: "Latvian",
		LanguageBCP47Code:   "lv",
		RegionDisplayName:   "Latvia",
		RegionBCP47Code:     "LV",
	},
	{
		Code:                "lt-LT",
		LanguageDisplayName: "Lithuanian",
		LanguageBCP47Code:   "lt",
		RegionDisplayName:   "Lithuania",
		RegionBCP47Code:     "LT",
	},
	{
		Code:                "ms-MY",
		LanguageDisplayName: "Malay",
		LanguageBCP47Code:   "ms",
		RegionDisplayName:   "Malaysia",
		RegionBCP47Code:     "MY",
	},
	{
		Code:                "ms-SG",
		LanguageDisplayName: "Malay",
		LanguageBCP47Code:   "ms",
		RegionDisplayName:   "Singapore",
		RegionBCP47Code:     "SG",
	},
	{
		Code:                "ml-IN",
		LanguageDisplayName: "Malayalam",
		LanguageBCP47Code:   "ml",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "mr-IN",
		LanguageDisplayName: "Marathi",
		LanguageBCP47Code:   "mr",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "no-NO",
		LanguageDisplayName: "Norwegian",
		LanguageBCP47Code:   "no",
		RegionDisplayName:   "Norway",
		RegionBCP47Code:     "NO",
	},
	{
		Code:                "pl-PL",
		LanguageDisplayName: "Polish",
		LanguageBCP47Code:   "pl",
		RegionDisplayName:   "Poland",
		RegionBCP47Code:     "PL",
	},
	{
		Code:                "pt-BR",
		LanguageDisplayName: "Portuguese",
		LanguageBCP47Code:   "pt",
		RegionDisplayName:   "Brazil",
		RegionBCP47Code:     "BR",
	},
	{
		Code:                "pt-PT",
		LanguageDisplayName: "Portuguese",
		LanguageBCP47Code:   "pt",
		RegionDisplayName:   "Portugal",
		RegionBCP47Code:     "PT",
	},
	{
		Code:                "ro-RO",
		LanguageDisplayName: "Romanian",
		LanguageBCP47Code:   "ro",
		RegionDisplayName:   "Romania",
		RegionBCP47Code:     "RO",
	},
	{
		Code:                "ru-BY",
		LanguageDisplayName: "Russian",
		LanguageBCP47Code:   "ru",
		RegionDisplayName:   "Belarus",
		RegionBCP47Code:     "BY",
	},
	{
		Code:                "ru-EE",
		LanguageDisplayName: "Russian",
		LanguageBCP47Code:   "ru",
		RegionDisplayName:   "Estonia",
		RegionBCP47Code:     "EE",
	},
	{
		Code:                "ru-IL",
		LanguageDisplayName: "Russian",
		LanguageBCP47Code:   "ru",
		RegionDisplayName:   "Israel",
		RegionBCP47Code:     "IL",
	},
	{
		Code:                "ru-RU",
		LanguageDisplayName: "Russian",
		LanguageBCP47Code:   "ru",
		RegionDisplayName:   "Russia",
		RegionBCP47Code:     "RU",
	},
	{
		Code:                "ru-UA",
		LanguageDisplayName: "Russian",
		LanguageBCP47Code:   "ru",
		RegionDisplayName:   "Ukraine",
		RegionBCP47Code:     "UA",
	},
	{
		Code:                "sr-Latn-RS",
		LanguageDisplayName: "Serbian",
		LanguageBCP47Code:   "sr",
		RegionDisplayName:   "Serbia",
		RegionBCP47Code:     "RS",
	},
	{
		Code:                "st-ZA",
		LanguageDisplayName: "Southern Sotho",
		LanguageBCP47Code:   "st",
		RegionDisplayName:   "South Africa",
		RegionBCP47Code:     "ZA",
	},
	{
		Code:                "sl-SI",
		LanguageDisplayName: "Slovenian",
		LanguageBCP47Code:   "sl",
		RegionDisplayName:   "Slovenia",
		RegionBCP47Code:     "SI",
	},
	{
		Code:                "es-AR",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Argentina",
		RegionBCP47Code:     "AR",
	},
	{
		Code:                "es-CL",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Chile",
		RegionBCP47Code:     "CL",
	},
	{
		Code:                "es-CO",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Colombia",
		RegionBCP47Code:     "CO",
	},
	{
		Code:                "es-CR",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Costa Rica",
		RegionBCP47Code:     "CR",
	},
	{
		Code:                "es-MX",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Mexico",
		RegionBCP47Code:     "MX",
	},
	{
		Code:                "es-PE",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Peru",
		RegionBCP47Code:     "PE",
	},
	{
		Code:                "es-ES",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "Spain",
		RegionBCP47Code:     "ES",
	},
	{
		Code:                "es-US",
		LanguageDisplayName: "Spanish",
		LanguageBCP47Code:   "es",
		RegionDisplayName:   "United States",
		RegionBCP47Code:     "US",
	},
	{
		Code:                "sw-KE",
		LanguageDisplayName: "Swahili",
		LanguageBCP47Code:   "sw",
		RegionDisplayName:   "Kenya",
		RegionBCP47Code:     "KE",
	},
	{
		Code:                "sv-SE",
		LanguageDisplayName: "Swedish",
		LanguageBCP47Code:   "sv",
		RegionDisplayName:   "Sweden",
		RegionBCP47Code:     "SE",
	},
	{
		Code:                "ta-IN",
		LanguageDisplayName: "Tamil",
		LanguageBCP47Code:   "ta",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "ta-MY",
		LanguageDisplayName: "Tamil",
		LanguageBCP47Code:   "ta",
		RegionDisplayName:   "Malaysia",
		RegionBCP47Code:     "MY",
	},
	{
		Code:                "ta-SG",
		LanguageDisplayName: "Tamil",
		LanguageBCP47Code:   "ta",
		RegionDisplayName:   "Singapore",
		RegionBCP47Code:     "SG",
	},
	{
		Code:                "te-IN",
		LanguageDisplayName: "Telugu",
		LanguageBCP47Code:   "te",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "th-TH",
		LanguageDisplayName: "Thai",
		LanguageBCP47Code:   "th",
		RegionDisplayName:   "Thailand",
		RegionBCP47Code:     "TH",
	},
	{
		Code:                "tr-TR",
		LanguageDisplayName: "Turkish",
		LanguageBCP47Code:   "tr",
		RegionDisplayName:   "Turkey",
		RegionBCP47Code:     "TR",
	},
	{
		Code:                "uk-UA",
		LanguageDisplayName: "Ukrainian",
		LanguageBCP47Code:   "uk",
		RegionDisplayName:   "Ukraine",
		RegionBCP47Code:     "UA",
	},
	{
		Code:                "ur-IN",
		LanguageDisplayName: "Urdu",
		LanguageBCP47Code:   "ur",
		RegionDisplayName:   "India",
		RegionBCP47Code:     "IN",
	},
	{
		Code:                "xh-ZA",
		LanguageDisplayName: "Xhosa",
		LanguageBCP47Code:   "xh",
		RegionDisplayName:   "South Africa",
		RegionBCP47Code:     "ZA",
	},
	{
		Code:                "zu-ZA",
		LanguageDisplayName: "Zulu",
		LanguageBCP47Code:   "zu",
		RegionDisplayName:   "South Africa",
		RegionBCP47Code:     "ZA",
	},
}

func TestLocales_GetSupportedLocales(t *testing.T) {
	t.Run("returns the full list of supported time zones", func(t *testing.T) {
		output := locales.GetSupportedLocales()
		require.EqualValues(t, expectedLocales, output)
	})
}

func TestLocales_LoadSupportedTag(t *testing.T) {
	t.Run("returns the language.Tag for the specified BCP 47 code", func(t *testing.T) {
		output, err := locales.LoadSupportedTag("en-US")
		require.NoError(t, err)
		require.Equal(t, language.AmericanEnglish, output)
	})

	t.Run("returns the language.Tag for the specified case-insensitive BCP 47 code", func(t *testing.T) {
		output, err := locales.LoadSupportedTag("en-us")
		require.NoError(t, err)
		require.Equal(t, language.AmericanEnglish, output)
	})

	t.Run("returns an error for an unsupported BCP 47 code", func(t *testing.T) {
		output, err := locales.LoadSupportedTag("en-XX")
		require.ErrorIs(t, err, locales.ErrUnsupportedLocale)
		require.Empty(t, output)
	})
}
