package datetime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Original PCRE Duration Representation:
// ^(-|\+)?P(?!$)((-|\+)?\d+(?:(\.|,)\d+)?Y)?((-|\+)?\d+(?:(\.|,)\d+)?M)?((-|\+)?\d+(?:(\.|,)\d+)?W)?((-|\+)?\d+(?:(\.|,)\d+)?D)?(T(?=(-|\+)?\d)((-|\+)?\d+(?:(\.|,)\d+)?H)?((-|\+)?\d+(?:(\.|,)\d+)?M)?((-|\+)?\d+(?:(\.|,)\d+)?S)?)?$
//          ^^^^^                                                                                                                  ^^^^^^^^^^^^^
//  Verify it is not just `P`                                                                                 Verify something else comes after T if T is included
//
// We can achieve this same representation in Golang's RE2 regex by replacing the negative & positive lookaheads with an additional regex expression check

// posneg contains regex for an optional pos or neg
const posneg = `(-|\+)?`

// num contains regex that represents a positive or negative number, possibly a decimal
var num = fmt.Sprintf(`%s\d+(?:(\.|,)\d+)?`, posneg)

// iso8601DurationPattern is derived from the perl representation in https://github.com/Urigo/graphql-scalars/blob/master/src/scalars/iso-date/Duration.ts
var iso8601DurationPattern = fmt.Sprintf(
	`^%sP(%sY)?(%sM)?(%sW)?(%sD)?(T(%sH)?(%sM)?(%sS)?)?$`,
	posneg, num, num, num, num, num, num, num,
)

// iso8601DurationRegExp is a regular expression for validating ISO 8601 duration strings
var iso8601DurationRegExp = regexp.MustCompile(iso8601DurationPattern)

// expectedMoreValuesRegExp replaces the two lookaheads that prevent just `P` and trailing `T` but no time component such as `P20DT`
var expectedMoreValuesRegExp = regexp.MustCompile(`(^P$)|(T$)`)

// Duration is an ISO8601 duration string
type Duration string

func ValidateDuration(dur string) bool {
	return iso8601DurationRegExp.MatchString(dur) && !expectedMoreValuesRegExp.MatchString(dur)
}

func ValidatePositiveDuration(dur string) bool {
	return ValidateDuration(dur) && !strings.Contains(dur, "-")
}

// UnmarshalJSON implements the graphql.Unmarshaler interface
func (d *Duration) UnmarshalJSON(data []byte) error {
	var rawDuration string
	err := json.Unmarshal(data, &rawDuration)

	if err != nil {
		return fmt.Errorf("value is not a valid ISO Duration, it must be a string: %w", err)
	}

	duration := strings.ToUpper(rawDuration)
	valid := ValidateDuration(duration)

	if !valid {
		return fmt.Errorf("value is not a valid ISO Duration: %s", rawDuration)
	}

	*d = Duration(duration)
	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Duration) UnmarshalGQL(v interface{}) error {
	rawDuration, ok := v.(string)
	if !ok {
		return errors.New("value is not a valid ISO Duration, it must be a string")
	}

	duration := strings.ToUpper(rawDuration)
	valid := ValidateDuration(duration)

	if !valid {
		return fmt.Errorf("value is not a valid ISO Duration: %s", rawDuration)
	}

	*d = Duration(duration)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Duration) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(string(d)))
}
