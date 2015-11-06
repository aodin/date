package date

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ISO8601Date uses ISO 8601 as a default for parsing and rendering
const ISO8601Date = "2006-01-02"

type Date struct{ time.Time }

func (date Date) format() string {
	return date.Time.Format(ISO8601Date)
}

// AddDays adds the given number of days to the date
func (date Date) AddDays(days int) Date {
	return Date{Time: date.Time.AddDate(0, 0, days)}
}

func (date Date) After(other Date) bool {
	return date.Time.After(other.Time)
}

func (date Date) Before(other Date) bool {
	return date.Time.Before(other.Time)
}

// String returns the Date as a string
func (date Date) String() string {
	return date.format()
}

// Equals returns true if the dates are equal
func (date Date) Equals(other Date) bool {
	return date.Time.Equal(other.Time)
}

// UnmarshalJSON converts a byte array into a Date
func (d *Date) UnmarshalJSON(text []byte) error {
	if string(text) == "null" {
		return nil
	}
	b := bytes.NewBuffer(text)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	value, err := time.Parse(ISO8601Date, s)
	if err != nil {
		return err
	}
	d.Time = value
	return nil
}

// MarshalJSON returns the JSON output of a Date.
// Null will return a zero value date.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.format() + `"`), nil
}

// Scan converts an SQL value into a Date
func (date *Date) Scan(value interface{}) error {
	date.Time = value.(time.Time)
	return nil
}

// Value returns the date formatted for insert into PostgreSQL
func (date Date) Value() (driver.Value, error) {
	return date.format(), nil
}

// Within returns true if the Date is within the range - inclusive
func (date Date) Within(term DateRange) bool {
	return !(date.Before(term.Start) || date.After(term.End))
}

// Today converts the local time to a Date
func Today() Date {
	return FromTime(time.Now())
}

// FromTime creates a Date from a time.Time
func FromTime(t time.Time) Date {
	return New(t.Year(), t.Month(), t.Day())
}

// New creates a new Date
func New(year int, month time.Month, day int) Date {
	// Remove all second and nano second information and mark as UTC
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// Parse converts a ISO 8601 date string to a Date, possibly returning an error
func Parse(value string) (Date, error) {
	return ParseUsingLayout(ISO8601Date, value)
}

// ParseUsingLayout calls Parse with a different date layout
func ParseUsingLayout(format, value string) (Date, error) {
	t, err := time.Parse(format, value)
	if err != nil {
		return Date{}, err
	}
	return Date{Time: t}, nil
}
