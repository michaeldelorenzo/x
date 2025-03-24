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

const DateFormat = "2006-01-02"

// Date is a wrapper of time.Time that only uses the YYYY-MM-DD when serializing. It represents an ISO 8601 date.
// Its internal time.Time's locale should always be set to UTC to ensure proper data comparison.
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}

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

func (d Date) String() string {
	return d.Time.Format(DateFormat)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Date) UnmarshalGQL(v interface{}) error {
	dateStr, ok := v.(string)
	if !ok {
		return errors.New("value is not a valid ISO Date, it must be a string")
	}

	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse date string into time: %w", err)
	}

	if d == nil {
		return fmt.Errorf("cannot scan %s into nil Date", parsed)
	}

	*d = TtoDate(parsed)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Date) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(d.String()))
}

// Value implements the driver.Valuer interface
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface
func (d *Date) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	parsed, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time but recieved %T when scanning into Date", src)
	}
	val := TtoDate(parsed)

	if d == nil {
		return fmt.Errorf("cannot scan %s into nil Date", val)
	}

	*d = val

	return nil
}

func (d *Date) UnmarshalText(data []byte) error {
	parsed, err := time.Parse(DateFormat, string(data))
	if err != nil {
		return err
	}

	*d = TtoDate(parsed)
	return nil
}

// coerceToUTC is an internal helper method that converts the internal date to UTC which should be used
// before performing any unsafe hour addition or subtraction
func (d Date) coerceToUTC() Date {
	if d.Location() != time.UTC {
		return TtoDate(d.Time)
	} else {
		return d
	}
}

// DaysSince computes the number of days since a previous date
func (d Date) DaysSince(previousDate Date) int {
	duration := d.coerceToUTC().Sub(previousDate.coerceToUTC().Time)
	return int(duration.Hours()) / 24
}

// DaysUntil computes the number of days until an upcoming date
func (d Date) DaysUntil(upcomingDate Date) int {
	duration := upcomingDate.coerceToUTC().Sub(d.coerceToUTC().Time)
	return int(duration.Hours()) / 24
}

// AddDays adds days to an existing date
func (d Date) AddDays(days int) Date {
	return d.AddDate(0, 0, days)
}

// AddDate returns the date corresponding to adding the given number of years, months, and days to d. For example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.
func (d Date) AddDate(years, months, days int) Date {
	return TtoDate(d.coerceToUTC().Time.AddDate(years, months, days))
}

// AtoDate provides a utility for scanning a string into a Date.
// If any error is encountered when parsing the string, it is swallowed and an empty Date is returned.
func AtoDate(input string) Date {
	parsed, _ := time.Parse(DateFormat, input)
	return TtoDate(parsed)
}

// TtoDate safely converts a time.Time into a Date, dropping all units below day and setting the tz to UTC
// for safe comparison between dates
func TtoDate(input time.Time) Date {
	return Date{Time: time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, time.UTC)}
}

// Today returns a new Date set to the current day
func Today() Date {
	return TtoDate(time.Now())
}
