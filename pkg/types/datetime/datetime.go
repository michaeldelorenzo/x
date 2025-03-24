package datetime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05.000Z07:00"

// DateTime is a wrapper of Time that always serializes with millisecond-level precision. It represents a full ISO 8601 timestamp with date and time.
type DateTime struct {
	time.Time
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateTimeFormat))
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("cannot unmarshal %s into nil DateTime", dateStr)
	}

	d.Time = parsed
	return nil
}

func (d DateTime) String() string {
	return d.Time.Format(DateTimeFormat)
}

// Value implements the driver.Valuer interface
func (d DateTime) Value() (driver.Value, error) {
	// format/parse to ensure we save with the same precision as we receive
	return time.Parse(DateTimeFormat, d.Time.Format(DateTimeFormat))
}

// Scan implements the sql.Scanner interface
func (d *DateTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	parsed, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time but recieved %T when scanning into DateTime", src)
	}

	val := DateTime{Time: parsed}

	if d == nil {
		return fmt.Errorf("cannot scan %s into nil DateTime", val.String())
	}

	*d = val

	return nil
}

// AtoDateTime provides a utility for scanning a string into a DateTime.
// If any error is encountered when parsing the string, it is swallowed and an empty DateTime is returned.
func AtoDateTime(input string) DateTime {
	parsed, _ := time.Parse(DateTimeFormat, input)
	return DateTime{Time: parsed}
}
