// Package locales contains the platform's source-of-truth list of supported locales
package locales

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldelorenzo/x/pkg/utils/sequence"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var ErrUnsupportedLocale = errors.New("provided locale is not supported")

type Locale struct {
	// The full locale code (i.e. es-US)
	Code string `json:"code" validate:"required" example:"es-US"`
	// The display name of the locale's base language (i.e. Spanish)
	LanguageDisplayName string `json:"language_display_name" validate:"required" example:"Spanish"`
	// The BCP 47 representation of the locale's base langauge (i.e. es)
	LanguageBCP47Code string `json:"language_bcp47_code" validate:"required" example:"es"`
	// The display name of the locale's region (i.e. United States)
	RegionDisplayName string `json:"region_display_name" validate:"required" example:"United States"`
	// The BCP 47 representation of the locale's region (i.e. US)
	RegionBCP47Code string `json:"region_bcp47_code" validate:"required" example:"US"`
}

func NewFromTag(tag language.Tag) *Locale {
	return &Locale{
		Code:                GetLocaleCode(tag),
		LanguageDisplayName: GetLanguageDisplayName(tag),
		LanguageBCP47Code:   GetLanguageCode(tag),
		RegionDisplayName:   GetRegionDisplayName(tag),
		RegionBCP47Code:     GetRegionCode(tag),
	}
}

func mustBuildTag(lang language.Tag, region language.Region) language.Tag {
	tag, err := language.Compose(lang, region)
	if err != nil {
		panic(err)
	}
	return tag
}

// Validate confirms that a string is a supported BCP 47 (case sensitive) language code
func Validate(code string) bool {
	for _, tag := range supportedLocales {
		localeCode := GetLocaleCode(tag)
		if localeCode == code {
			return true
		}
	}

	return false
}

// GetLocaleCode returns the unique code comprised of `<ISO 639-1 Language Code>-<ISO 3166-1 Region Code>`
// the codes should match Xcode and [this list](https://www.oracle.com/java/technologies/javase/jdk8-jre8-suported-locales.html)
func GetLocaleCode(tag language.Tag) string {
	return tag.String()
}

// GetLanguageCode returns the BCP 47 representation of the locale's base language
func GetLanguageCode(tag language.Tag) string {
	language, _ := tag.Base()
	return language.String()
}

// GetRegionCode returns the BCP 47 representation of the locale's region
func GetRegionCode(tag language.Tag) string {
	region, _ := tag.Region()
	return region.String()
}

// languageDisplayNameOverrides overrides the default English naming of the language
var languageDisplayNameOverrides = map[string]string{
	"no-NO":  "Norwegian", // display name should be `Norwegian` but is `Norwegian Bokmål`
	"luo-KE": "Dholuo",    // display name should be `Dholuo` but is `Luo`
}

// GetLanguageDisplayName returns the display name of a tag's base language
func GetLanguageDisplayName(tag language.Tag) string {
	if name, ok := languageDisplayNameOverrides[tag.String()]; ok {
		return name
	}
	language, _ := tag.Base()
	return display.AmericanEnglish.Languages().Name(language)
}

// regionDisplayNameOverrides overrides the default English naming of the region
var regionDisplayNameOverrides = map[string]string{
	// Czech Republic renamed themselves in 2016 to Czechia,
	// if we want to use the original name just uncomment the below override
	// "CZ": "Czech Republic",
}

// GetRegionDisplayName returns the display name of a tag's region
func GetRegionDisplayName(tag language.Tag) string {
	region, _ := tag.Region()

	if name, ok := regionDisplayNameOverrides[region.String()]; ok {
		return name
	}

	return display.AmericanEnglish.Regions().Name(region)
}

// GetSupportedLanguageTags returns a list of supported language tags
func GetSupportedLanguageTags() []language.Tag {
	return supportedLocales
}

// GetSupportedLocales returns a list of supported locales
func GetSupportedLocales() []Locale {
	return sequence.Map(GetSupportedLanguageTags(), func(tag language.Tag) Locale {
		return *NewFromTag(tag)
	})
}

// LoadSupportedTag returns the loaded location for a case insensitive BCP 47 code or an error if it is not found in the supported list
func LoadSupportedTag(code string) (language.Tag, error) {
	tag, ok := supportedLowerCaseLocaleMap[strings.ToLower(code)]
	if !ok {
		return tag, fmt.Errorf("failed to load language tag for %s: %w", code, ErrUnsupportedLocale)
	}

	return tag, nil
}

// regions
var argentina = language.MustParseRegion("ar")
var australia = language.MustParseRegion("au")
var austria = language.MustParseRegion("at")
var belarus = language.MustParseRegion("by")
var belgium = language.MustParseRegion("be")
var brazil = language.MustParseRegion("br")
var bulgaria = language.MustParseRegion("bg")
var canada = language.MustParseRegion("ca")
var chile = language.MustParseRegion("cl")
var china = language.MustParseRegion("cn")
var colombia = language.MustParseRegion("co")
var costaRica = language.MustParseRegion("cr")
var croatia = language.MustParseRegion("hr")
var czechia = language.MustParseRegion("cz") // Czech Republic renamed themselves in 2016 to Czechia
var denmark = language.MustParseRegion("dk")
var egypt = language.MustParseRegion("eg")
var estonia = language.MustParseRegion("ee")

var finland = language.MustParseRegion("fi")
var france = language.MustParseRegion("fr")
var georgia = language.MustParseRegion("ge")
var germany = language.MustParseRegion("de")
var greece = language.MustParseRegion("gr")
var hungary = language.MustParseRegion("hu")
var india = language.MustParseRegion("in")
var ireland = language.MustParseRegion("ie")
var israel = language.MustParseRegion("il")
var italy = language.MustParseRegion("it")
var japan = language.MustParseRegion("jp")
var kenya = language.MustParseRegion("ke")
var latvia = language.MustParseRegion("lv")
var lebanon = language.MustParseRegion("lb")
var lithuania = language.MustParseRegion("lt")
var malaysia = language.MustParseRegion("my")
var mexico = language.MustParseRegion("mx")
var netherlands = language.MustParseRegion("nl")
var norway = language.MustParseRegion("no")
var peru = language.MustParseRegion("pe")
var poland = language.MustParseRegion("pl")
var portugal = language.MustParseRegion("pt")
var qatar = language.MustParseRegion("qa")
var romania = language.MustParseRegion("ro")
var russia = language.MustParseRegion("ru")
var serbia = language.MustParseRegion("rs")
var singapore = language.MustParseRegion("sg")
var slovenia = language.MustParseRegion("si")
var southAfrica = language.MustParseRegion("za")
var southKorea = language.MustParseRegion("kr")
var spain = language.MustParseRegion("es")
var sweden = language.MustParseRegion("se")
var switzerland = language.MustParseRegion("ch")
var thailand = language.MustParseRegion("th")
var turkey = language.MustParseRegion("tr")
var ukraine = language.MustParseRegion("ua")
var unitedArabEmirates = language.MustParseRegion("ae")
var unitedKingdom = language.MustParseRegion("gb")
var unitedStates = language.MustParseRegion("us")

// missing languages
var belarusian = language.MustParse("be")

// var dholuo = language.MustParse("luo")
var irish = language.MustParse("ga")
var norwegian = language.MustParse("no")

// var norwegianBokmål = language.MustParse("nb")
// var norwegianNyorsk = language.MustParse("nn")
var sesotho = language.MustParse("st")
var xhosa = language.MustParse("xh")

// Commented out languages below were originally added because the legacy platform supported them,
// but were then removed because it was determined they were not yet fully supported in our Mobile App.
// They have been left commented out for ease of reimplementation when/if they are readded in VNext.

// supportedLocales contains the list of supported locales in the platform
var supportedLocales = []language.Tag{
	// en-US and en-GB languages listed first for ease of use
	mustBuildTag(language.AmericanEnglish, unitedStates),
	mustBuildTag(language.BritishEnglish, unitedKingdom),

	mustBuildTag(language.Afrikaans, southAfrica),
	mustBuildTag(language.Arabic, egypt),
	mustBuildTag(language.Arabic, israel),
	mustBuildTag(language.Arabic, lebanon),
	mustBuildTag(language.Arabic, qatar),
	mustBuildTag(language.Arabic, unitedArabEmirates),
	mustBuildTag( /*   */ belarusian, belarus),
	mustBuildTag(language.Bengali, india),
	mustBuildTag(language.Bulgarian, bulgaria),
	mustBuildTag(language.SimplifiedChinese, china),
	mustBuildTag(language.SimplifiedChinese, malaysia),
	mustBuildTag(language.SimplifiedChinese, singapore),
	mustBuildTag(language.Croatian, croatia),
	mustBuildTag(language.Czech, czechia),
	mustBuildTag(language.Danish, denmark),
	// mustBuildTag( /*   */ dholuo, kenya),
	mustBuildTag(language.Dutch, belgium),
	mustBuildTag(language.Dutch, netherlands),
	mustBuildTag(language.English, australia),
	mustBuildTag(language.English, canada),
	mustBuildTag(language.English, india),
	mustBuildTag(language.English, ireland),
	mustBuildTag(language.English, kenya),
	mustBuildTag(language.English, malaysia),
	mustBuildTag(language.English, singapore),
	mustBuildTag(language.English, southAfrica),
	mustBuildTag(language.Estonian, estonia),
	mustBuildTag(language.Finnish, finland),
	mustBuildTag(language.Swedish, finland),
	mustBuildTag(language.French, belgium),
	mustBuildTag(language.French, canada),
	mustBuildTag(language.French, france),
	mustBuildTag(language.French, switzerland),
	mustBuildTag(language.Georgian, georgia),
	mustBuildTag(language.German, austria),
	mustBuildTag(language.German, belgium),
	mustBuildTag(language.German, germany),
	mustBuildTag(language.German, switzerland),
	mustBuildTag(language.Greek, greece),
	mustBuildTag(language.Gujarati, india),
	mustBuildTag(language.Hebrew, israel),
	mustBuildTag(language.Hindi, india),
	mustBuildTag(language.Hungarian, hungary),
	mustBuildTag( /*   */ irish, ireland),
	mustBuildTag(language.Italian, italy),
	mustBuildTag(language.Italian, switzerland),
	mustBuildTag(language.Japanese, japan),
	mustBuildTag(language.Kannada, india),
	mustBuildTag(language.Korean, southKorea),
	mustBuildTag(language.Latvian, latvia),
	mustBuildTag(language.Lithuanian, lithuania),
	mustBuildTag(language.Malay, malaysia),
	mustBuildTag(language.Malay, singapore),
	mustBuildTag(language.Malayalam, india),
	mustBuildTag(language.Marathi, india),
	mustBuildTag( /*   */ norwegian, norway),
	// mustBuildTag( /*   */ norwegianBokmål, norway),
	// mustBuildTag( /*   */ norwegianNyorsk, norway),
	mustBuildTag(language.Polish, poland),
	mustBuildTag(language.Portuguese, brazil),
	mustBuildTag(language.Portuguese, portugal),
	mustBuildTag(language.Romanian, romania),
	mustBuildTag(language.Russian, belarus),
	mustBuildTag(language.Russian, estonia),
	mustBuildTag(language.Russian, israel),
	// mustBuildTag(language.Russian, latvia),
	// mustBuildTag(language.Russian, lithuania),
	mustBuildTag(language.Russian, russia),
	mustBuildTag(language.Russian, ukraine),
	mustBuildTag(language.SerbianLatin, serbia),
	mustBuildTag( /*   */ sesotho, southAfrica),
	mustBuildTag(language.Slovenian, slovenia),
	mustBuildTag(language.Spanish, argentina),
	mustBuildTag(language.Spanish, chile),
	mustBuildTag(language.Spanish, colombia),
	mustBuildTag(language.Spanish, costaRica),
	mustBuildTag(language.Spanish, mexico),
	mustBuildTag(language.Spanish, peru),
	mustBuildTag(language.Spanish, spain),
	mustBuildTag(language.Spanish, unitedStates),
	mustBuildTag(language.Swahili, kenya),
	mustBuildTag(language.Swedish, sweden),
	mustBuildTag(language.Tamil, india),
	mustBuildTag(language.Tamil, malaysia),
	mustBuildTag(language.Tamil, singapore),
	mustBuildTag(language.Telugu, india),
	mustBuildTag(language.Thai, thailand),
	mustBuildTag(language.Turkish, turkey),
	mustBuildTag(language.Ukrainian, ukraine),
	mustBuildTag(language.Urdu, india),
	mustBuildTag( /*   */ xhosa, southAfrica),
	mustBuildTag(language.Zulu, southAfrica),
}

// supportedLowerCaseLocaleMap contains a mapping of lower-cased time zones to their time.Location struct counterparts
var supportedLowerCaseLocaleMap = sequence.Reduce(supportedLocales, map[string]language.Tag{}, func(acc map[string]language.Tag, el language.Tag) map[string]language.Tag {
	acc[strings.ToLower(GetLocaleCode(el))] = el
	return acc
})
