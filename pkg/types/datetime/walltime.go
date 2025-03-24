package datetime

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

const TimeFormat = "15:04:05"

// WallTime is a wrapper of time that only outputs the HH:mm:ss during serialization. It represents an ISO 8601 time.
type WallTime struct {
	time.Time
}

func (t WallTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(TimeFormat))
}

func (t *WallTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}

	parsed, err := time.Parse(TimeFormat, dateStr)
	if err != nil {
		return err
	}

	if t == nil {
		return fmt.Errorf("cannot unmarshal %s into nil WallTime", dateStr)
	}

	t.Time = parsed
	return nil
}

func (t WallTime) String() string {
	return t.Time.Format(TimeFormat)
}

// Value implements the driver.Valuer interface
func (t WallTime) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan implements the sql.Scanner interface
func (t *WallTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	parsed, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time but recieved %T when scanning into WallTime", src)
	}
	val := WallTime{Time: parsed}

	if t == nil {
		return fmt.Errorf("cannot scan %s into nil WallTime", val)
	}

	*t = val

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *WallTime) UnmarshalGQL(v interface{}) error {
	dateStr, ok := v.(string)
	if !ok {
		return errors.New("value is not a valid ISO Time, it must be a string")
	}

	parsed, err := time.Parse(TimeFormat, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse time string into wall time: %w", err)
	}

	if t == nil {
		return fmt.Errorf("cannot scan %s into nil WallTime", parsed)
	}

	t.Time = parsed
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t WallTime) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(t.String()))
}

func (d *WallTime) UnmarshalText(data []byte) error {
	parsed, err := time.Parse(TimeFormat, string(data))
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

// AtoWallTime provides a utility for scanning a string into a WallTime.
// If any error is encountered when parsing the string, it is swallowed and an empty WallTime is returned.
func AtoWallTime(input string) WallTime {
	parsed, _ := time.Parse(TimeFormat, input)
	return WallTime{Time: parsed}
}
