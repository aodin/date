package date

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DateRange is a start and end date
type DateRange struct {
	Start   Date `json:"start"`
	End     Date `json:"end"`
	isEmpty bool
}

// Contains returns true if the given range is entirely within the
// the range - inclusive
func (term DateRange) Contains(other DateRange) bool {
	return !term.DoesNotContain(other)
}

func (term DateRange) DoesNotContain(other DateRange) bool {
	return other.Start.Before(term.Start) || other.End.After(term.End)
}

// Equals returns true if the term has equal start and end dates
func (term DateRange) Equals(other DateRange) bool {
	return term.isEmpty == other.isEmpty && term.Start.Equals(other.Start) && term.End.Equals(other.End)
}

// Error returns an error if there is both a start and end date and the given
// start date is not before the end date.
func (term DateRange) Error() error {
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
func (term DateRange) IsInfinity() bool {
	return term.IsZero()
}

// IsEmpty returns true if the term is zero and has isEmpty = true
func (term DateRange) IsEmpty() bool {
	return term.isEmpty && term.IsZero()
}

// IsZero returns true if the start and end dates are both zero
func (term DateRange) IsZero() bool {
	return term.Start.IsZero() && term.End.IsZero()
}

func isEmptyDateRange(value string) bool {
	return strings.ToLower(value) == "empty"
}

// splitDateRange divides a term into start and end date strings
func splitDateRange(value string) (string, string, error) {
	p := strings.SplitN(value, ",", 2)
	if len(p) != 2 || p[0] == "" || p[1] == "" {
		return "", "", fmt.Errorf("date: failed to parse date range '%s'", value)
	}
	return strings.ToLower(p[0][1:]), strings.ToLower(p[1][:len(p[1])-1]), nil
}

// Scan converts the given database value to a DateRange,
// possibly returning an error if the conversion failed
func (term *DateRange) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("date: failed to convert date range to []byte")
	}

	// Zero ranges return "empty"
	if isEmptyDateRange(string(b)) {
		term.isEmpty = true
		return nil
	}

	// Otherwise, parse the given SQL date range
	start, end, err := splitDateRange(string(b))
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
func (term DateRange) String() string {
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

func (term DateRange) Intersection(other DateRange) (intersect DateRange) {
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

// MarshalJSON returns the JSON output of a DateRange.
// Empty ranges will return null
func (term DateRange) MarshalJSON() ([]byte, error) {
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

// Union creates the union of two DateRange types. If there is a gap
// between the two range it is included.
func (term DateRange) Union(other DateRange) (union DateRange) {
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
func (term DateRange) Value() (driver.Value, error) {
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

// Empty creates an empty DateRange
func Empty() (term DateRange) {
	term.isEmpty = true
	return
}

// Forever creates a DateRange without a start or end date
func Forever() (term DateRange) {
	return
}

// Infinity is an alias for Forever
func Infinity() DateRange {
	return Forever()
}

// Never is an alias for Empty
func Never() DateRange {
	return Empty()
}

// Range creates a DateRange with the given start and end dates
func Range(start, end Date) (term DateRange) {
	term.Start = start
	term.End = end
	return
}

// EntireMonth creates a DateRange that includes the entire month
func EntireMonth(year int, month time.Month) DateRange {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1)
	return Range(New(first.Date()), New(last.Date()))
}

// EntireYear creates a DateRange that includes the entire year
func EntireYear(year int) DateRange {
	first := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(1, 0, -1)
	return Range(New(first.Date()), New(last.Date()))
}

func SingleDay(date Date) DateRange {
	return Range(date, date)
}

func OnlyToday() DateRange {
	return SingleDay(Today())
}

func StartBoundedRange(start Date) (term DateRange) {
	term.Start = start
	return
}
