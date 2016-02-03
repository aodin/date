package date

import (
	"encoding/json"
	"testing"
)

var (
	empty     = Empty()
	year2015  = EntireYear(2015)
	nov       = EntireMonth(2015, 11)
	dec       = EntireMonth(2015, 12)
	novOnward = Range{Start: nov.Start}
	decOnward = Range{Start: dec.Start}
	untilDec  = Range{End: nov.End}
	jan       = EntireMonth(2016, 1)
	feb       = EntireMonth(2016, 2)
)

func TestRange(t *testing.T) {
	if !Infinity().IsInfinity() {
		t.Error("Infinity should be infinity")
	}
	if !Empty().IsEmpty() {
		t.Error("Empty should be empty")
	}
	if !Never().Equals(Empty()) {
		t.Error("Never should equal Empty")
	}
	if !Forever().Equals(Infinity()) {
		t.Error("Forever should equal Infinity")
	}

	if Forever().Equals(Empty()) {
		t.Error("Forever should not equal Empty")
	}
}

func TestRange_Contains(t *testing.T) {
	if !year2015.Contains(dec) {
		t.Error("year 2015 should contain December 2015")
	}
	if !dec.DoesNotContain(year2015) {
		t.Error("December 2015 should not contain the year 2015")
	}

	if !novOnward.Contains(dec) {
		t.Error("November 2015 onward should contain December 2015")
	}
	if year2015.Contains(novOnward) {
		t.Error("Year 2015 should not contain November 2015 onward")
	}

	if !novOnward.Contains(novOnward) {
		t.Error("November 2015 onward should contain itself")
	}

	if !novOnward.Contains(decOnward) {
		t.Error("November 2015 onward should contain December 2015 onward")
	}
	if decOnward.Contains(novOnward) {
		t.Error("December 2015 onward should not contain November 2015 onward")
	}
}

var daysTests = []struct {
	want, have int
}{
	{0, Empty().Days()},
	{1, OnlyToday().Days()},
	{365, EntireYear(2015).Days()},
	{366, EntireYear(2016).Days()},
	{29, EntireMonth(2016, 2).Days()},
}

func TestRange_Days(t *testing.T) {
	for _, test := range daysTests {
		if test.want != test.have {
			t.Errorf("Range Days() want=%d have=%d", test.want, test.have)
		}
	}
}

func TestRange_Error(t *testing.T) {
	if Never().Error() != nil {
		t.Error("Never should not error")
	}

	// Unbounded ranges are allowed
	end := Range{End: New(2015, 3, 1)}
	if end.Error() != nil {
		t.Error("Unbounded start dates should not error")
	}

	start := Range{Start: New(2015, 3, 1)}
	if start.Error() != nil {
		t.Error("Unbounded end dates should not error")
	}

	var invalid Range
	invalid.Start = New(2015, 3, 2)
	invalid.End = New(2015, 3, 1)
	if invalid.Error() == nil {
		t.Error("Invalid ranges should error")
	}
}

var intersectionTests = []struct {
	want, have Range
}{
	{dec, novOnward.Intersection(dec)},
	{nov, year2015.Intersection(nov)},
	{decOnward, decOnward.Intersection(novOnward)},
	{decOnward, novOnward.Intersection(decOnward)},
	{nov, novOnward.Intersection(untilDec)},
}

func TestRange_Intersection(t *testing.T) {
	if !empty.Intersection(nov).IsZero() {
		t.Error("An empty range should have a zero intersection")
	}

	for _, test := range intersectionTests {
		if test.want != test.have {
			t.Errorf(
				"Range Intersection() want=%v have=%v", test.want, test.have,
			)
		}
	}
}

func TestRange_Overlaps(t *testing.T) {
	if nov.Overlaps(dec) {
		t.Error("November should not overlap December")
	}
	if !dec.Overlaps(year2015) {
		t.Error("December 2015 should overlap the year 2015")
	}
	if !nov.Overlaps(SingleDay(New(2015, 11, 30))) {
		t.Error("November 2015 should overlap 2015-11-30")
	}
	if !novOnward.Overlaps(dec) {
		t.Error("November 2015 onward should overlap December 2015")
	}
}

func TestRange_Marshal(t *testing.T) {
	// Empty ranges should render as null
	b, err := json.Marshal(Never())
	if err != nil {
		t.Fatal("json.Marshal of Never() should not error")
	}
	if string(b) != "null" {
		t.Error(`json.Marshal of Never() should be null`)
	}

	// Infinite ranges should render as null start and end dates
	b, err = json.Marshal(Infinity())
	if err != nil {
		t.Fatal("json.Marshal of Infinity() should not error")
	}
	if string(b) != `{"start":null,"end":null}` {
		t.Error(
			`json.Marshal of Infinity() should be {"start":null,"end":null}`,
		)
	}
}

var stringTests = []struct {
	want, have string
}{
	{"never", Never().String()},
	{"forever", Forever().String()},
	{"2016-02-01 to 2016-02-29", EntireMonth(2016, 2).String()},
	{"until 2016-02-29", Range{End: New(2016, 2, 29)}.String()},
	{"2016-02-01 onward", Range{Start: New(2016, 2, 1)}.String()},
}

func TestRange_String(t *testing.T) {
	for _, test := range stringTests {
		if test.want != test.have {
			t.Errorf("Range String() want=%s have=%s", test.want, test.have)
		}
	}
}

func TestRange_Union(t *testing.T) {
	union := year2015.Union(jan)
	if New(2015, 1, 1) != union.Start {
		t.Error(
			"The union of 2015 and January 2016 should start on 2015-01-01",
		)
	}
	if New(2016, 1, 31) != union.End {
		t.Error(
			"The union of 2015 and January 2016 should end on 2016-01-31",
		)
	}

	union = jan.Union(feb)
	if jan.Start != union.Start {
		t.Error(
			"The union of January and February 2016 should start on 2016-01-01",
		)
	}
	if feb.End != union.End {
		t.Error(
			"The union of January and February 2016 should end on 2016-02-29",
		)
	}

	if decOnward.Union(novOnward) != novOnward {
		t.Error("The union of November onward and December onward should be Novermber onward")
	}
	if untilDec.Union(decOnward) != Forever() {
		t.Error("The union of until December and December onward should be forever")
	}

	if Empty().Union(Empty()) != Empty() {
		t.Error("The union of two empty ranges should be empty")
	}
	if Empty().Union(Forever()) != Forever() {
		t.Error("The union of a forever range should be forever")
	}
	if Forever().Union(Empty()) != Forever() {
		t.Error("The union of a forever range should be forever")
	}
}

func TestRange_Unmarshal(t *testing.T) {
	// Unmarshaling should overwrite values
	open := EntireMonth(2015, 2)

	raw := `{"start":"2015-03-01","end":null}`
	if json.Unmarshal([]byte(raw), &open) != nil {
		t.Error("json.Unmarshal should not error")
	}

	if New(2015, 3, 1) != open.Start {
		t.Error("The start date after json.Unmarshal should be 2015-03-01")
	}
	if !open.End.IsZero() {
		t.Error("The end date after json.Unmarshal should be zero")
	}

	// TODO nulls should be unmarshaled as empty ranges
	// raw = `null`
	// var zero Range
	// if json.Unmarshal([]byte(raw), &zero) != nil {
	// 	t.Error("json.Unmarshal should not error for null ranges")
	// }
	// if !zero.IsEmpty() {
	// 	t.Error("null ranges after json.Unmarshal should be empty")
	// }
}
