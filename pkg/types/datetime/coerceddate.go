package datetime

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CoercedDate is a wrapper of Date that  coerces date strings longer than YYYY-MM-DD into YYYY-MM-DD when serializing.
// It represents an ISO 8601 date.
type CoercedDate struct{ Date }

const dateFormatLength = len(DateFormat)

func (d *CoercedDate) UnmarshalJSON(data []byte) error {
	var inputStr string
	err := json.Unmarshal(data, &inputStr)
	if err != nil {
		return err
	}

	if len(inputStr) < dateFormatLength {
		return fmt.Errorf("value cannot be coerced into a date, it must be at least %d characters long", dateFormatLength)
	}
	dateStr := inputStr[0:dateFormatLength]

	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("cannot unmarshal %s into nil Date", dateStr)
	}

	d.Time = parsed
	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *CoercedDate) UnmarshalGQL(v interface{}) error {
	inputStr, ok := v.(string)
	if !ok {
		return errors.New("value is not a valid ISO Date, it must be a string")
	}

	if len(inputStr) < dateFormatLength {
		return fmt.Errorf("value cannot be coerced into Date, it must be at least %d characters long", dateFormatLength)
	}
	dateStr := inputStr[0:dateFormatLength]

	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse date string into time: %w", err)
	}

	if d == nil {
		return fmt.Errorf("cannot scan %s into nil Date", parsed)
	}

	d.Time = parsed
	return nil
}

func (d *CoercedDate) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	parsed, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time but recieved %T when scanning into Date", src)
	}
	val := CoercedDate{Date: Date{Time: parsed}}

	if d == nil {
		return fmt.Errorf("cannot scan %s into nil CoercedDate", val)
	}

	*d = val

	return nil
}

func (d *CoercedDate) UnmarshalText(data []byte) error {
	if len(string(data)) < dateFormatLength {
		return fmt.Errorf("value cannot be coerced into Date, it must be at least %d characters long", dateFormatLength)
	}

	parsed, err := time.Parse(DateFormat, string(data)[0:dateFormatLength])
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}
