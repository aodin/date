package date

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Range is a start and end date
type Range struct {
	Start   Date `json:"start"`
	End     Date `json:"end"`
	isEmpty bool
}

// Contains returns true if the given range is entirely within the
// the range - inclusive
func (term Range) Contains(other Range) bool {
	return !term.DoesNotContain(other)
}

func (term Range) DoesNotContain(other Range) bool {
	return other.Start.Before(term.Start) || other.End.After(term.End)
}

// Equals returns true if the term has equal start and end dates
func (term Range) Equals(other Range) bool {
	return term.isEmpty == other.isEmpty && term.Start.Equals(other.Start) && term.End.Equals(other.End)
}

// Error returns an error if there is both a start and end date and the given
// start date is not before the end date.
func (term Range) Error() error {
	if term.IsZero() {
		return nil
	}
	// One day only is allowed
	if term.Start.After(term.End) {
		return fmt.Errorf("Start date cannot be after the end date")
	}
	return nil
}

// IsInfinity is an alias for IsZero
func (term Range) IsInfinity() bool {
	return term.IsZero()
}

// IsEmpty returns true if the term is zero and has isEmpty = true
func (term Range) IsEmpty() bool {
	return term.isEmpty && term.IsZero()
}

// IsZero returns true if the start and end dates are both zero
func (term Range) IsZero() bool {
	return term.Start.IsZero() && term.End.IsZero()
}

func isEmptyRange(value string) bool {
	return strings.ToLower(value) == "empty"
}

// splitRange divides a term into start and end date strings
func splitRange(value string) (string, string, error) {
	p := strings.SplitN(value, ",", 2)
	if len(p) != 2 || p[0] == "" || p[1] == "" {
		return "", "", fmt.Errorf("date: failed to parse date range '%s'", value)
	}
	return strings.ToLower(p[0][1:]), strings.ToLower(p[1][:len(p[1])-1]), nil
}

// Scan converts the given database value to a Range,
// possibly returning an error if the conversion failed
func (term *Range) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("date: failed to convert date range to []byte")
	}

	// Zero ranges return "empty"
	if isEmptyRange(string(b)) {
		term.isEmpty = true
		return nil
	}

	// Otherwise, parse the given SQL date range
	start, end, err := splitRange(string(b))
	if err != nil {
		return err
	}

	if start == "infinity" || start == "" {
		// do nothing
	} else {
		var startDate Date
		if startDate, err = Parse(start); err != nil {
			return err
		}
		term.Start = startDate
	}

	if end == "infinity" || end == "" {
		return nil
	}

	var endDate Date
	if endDate, err = Parse(end); err != nil {
		return err
	}

	// Remove a single day from the date (it is exclusive - we want inclusive)
	endDate = endDate.AddDays(-1)
	term.End = endDate
	return nil
}

// String returns a string representation of the date range
func (term Range) String() string {
	if term.IsZero() {
		return "forever"
	}
	if term.Start.IsZero() {
		return fmt.Sprintf("until %s", term.End)
	}
	if term.End.IsZero() {
		return fmt.Sprintf("%s onward", term.Start)
	}
	return fmt.Sprintf("%s to %s", term.Start, term.End)
}

func (term Range) Intersection(other Range) (intersect Range) {
	// If either range is empty then the intersection is empty
	if term.IsEmpty() || other.IsEmpty() {
		intersect.isEmpty = true
		return
	}

	if other.Start.Within(term) {
		intersect.Start = other.Start
	} else if term.Start.Within(other) {
		intersect.End = term.Start
	} else {
		intersect.isEmpty = true
		return
	}

	if other.End.Within(term) {
		intersect.End = other.End
	} else if term.End.Within(other) {
		intersect.End = term.End
	} else {
		intersect.isEmpty = true
		return
	}
	return
}

// MarshalJSON returns the JSON output of a Range.
// Empty ranges will return null
func (term Range) MarshalJSON() ([]byte, error) {
	if term.IsEmpty() {
		return []byte("null"), nil
	}
	start, err := json.Marshal(term.Start)
	if err != nil {
		return nil, err
	}
	end, err := json.Marshal(term.End)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{"start":%s,"end":%s}`, start, end)), nil
}

// Union creates the union of two Range types. If there is a gap
// between the two range it is included.
func (term Range) Union(other Range) (union Range) {
	if term.Start.Before(other.Start) {
		union.Start = term.Start
	} else {
		union.Start = other.Start
	}
	if term.End.After(other.End) {
		union.End = term.End
	} else {
		union.End = other.End
	}
	return
}

// Value prepares the nullable term for the database
func (term Range) Value() (driver.Value, error) {
	if term.IsZero() {
		return "[,]", nil
	}
	if term.Start.IsZero() {
		return fmt.Sprintf("[,'%s']", term.End), nil
	}
	if term.End.IsZero() {
		return fmt.Sprintf("['%s',]", term.Start), nil
	}
	return fmt.Sprintf("['%s','%s']", term.Start, term.End), nil
}

// Empty creates an empty Range
func Empty() (term Range) {
	term.isEmpty = true
	return
}

// Forever creates a Range without a start or end date
func Forever() (term Range) {
	return
}

// Infinity is an alias for Forever
func Infinity() Range {
	return Forever()
}

// Never is an alias for Empty
func Never() Range {
	return Empty()
}

// NewRange creates a Range with the given start and end dates
func NewRange(start, end Date) (term Range) {
	term.Start = start
	term.End = end
	return
}

// EntireMonth creates a Range that includes the entire month
func EntireMonth(year int, month time.Month) (term Range) {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	term.End = FromTime(first.AddDate(0, 1, -1))
	term.Start = FromTime(first)
	return
}

// EntireYear creates a Range that includes the entire year
func EntireYear(year int) (term Range) {
	first := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	term.End = FromTime(first.AddDate(1, 0, -1))
	term.Start = FromTime(first)
	return
}

func SingleDay(date Date) Range {
	return NewRange(date, date)
}

func OnlyToday() Range {
	return SingleDay(Today())
}

func StartBoundedRange(start Date) (term Range) {
	term.Start = start
	return
}
