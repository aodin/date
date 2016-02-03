package date

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	day := New(2015, 3, 1)
	if day.String() != "2015-03-01" {
		t.Error("The string output of March 1st should be 2015-03-01")
	}

	output, err := json.Marshal(day)
	if err != nil {
		t.Error("JSON marshaling of dates should not error")
	}
	if string(output) != `"2015-03-01"` {
		t.Error(`JSON marshaling of March 1st should be "2015-03-01"`)
	}

	// Zero dates should return null
	var zero Date
	output, err = json.Marshal(zero)
	if err != nil {
		t.Error("JSON marshaling of zero dates should not error")
	}
	if string(output) != "null" {
		t.Error("json.Marshal of a zero date should be null")
	}

	nextDay := New(2015, 3, 2)
	if !nextDay.Equals(day.AddDays(1)) {
		t.Error("The day after March 1st should be March 2nd")
	}

	parsed, err := Parse("2015-03-01")
	if err != nil {
		t.Error("Parsing of properly formatted dates should not error")
	}
	if !parsed.Equals(day) {
		t.Error("The parsed string should equal March 1st")
	}

	if day.UnmarshalJSON([]byte(`"2015-03-01"`)) != nil {
		t.Error("UnmarshalJSON of a valid slice of bytes should not error")
	}

	// Parsing null should return a zero date
	if zero.UnmarshalJSON([]byte("null")) != nil {
		t.Error("Unmarshaling a null date should not error")
	}
	if !zero.IsZero() {
		t.Error("A null date should unmarshal to zero")
	}

	jan1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	if !jan1.Before(day.Time) {
		t.Error("January 1st should be before March 1st")
	}
}

func TestDate_Within(t *testing.T) {
	march1 := New(2015, 3, 1)
	dec1 := New(2015, 12, 1)

	feb := EntireMonth(2015, 2)
	march := EntireMonth(2015, 3)

	if march1.Within(feb) {
		t.Error("March 1st should not be within February")
	}
	if march1 != march.Start {
		t.Error("March 1st should equal the start of March")
	}
	if !march1.Within(march) {
		t.Error("March 1st should be within March")
	}

	// Test unbounded ranges
	novOnward := Range{Start: New(2015, 11, 1)}
	beforeNov := Range{End: New(2015, 10, 31)}

	if !dec1.Within(novOnward) {
		t.Error("December 1st should be within November onward")
	}
	if dec1.Within(beforeNov) {
		t.Error("December 1st should not be within before November")
	}

	if !march1.Within(beforeNov) {
		t.Error("March 1st should be within before November")
	}
	if march1.Within(novOnward) {
		t.Error("March 1st should not be within November onward")
	}
}
