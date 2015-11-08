package date

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	assert.True(t, Infinity().IsInfinity())
	assert.True(t, Empty().IsEmpty())

	assert.True(t, Never().Equals(Empty()))
	assert.True(t, Forever().Equals(Infinity()))
	assert.False(t, Forever().Equals(Empty()))
}

func TestRange_Contains(t *testing.T) {
	year2015 := EntireYear(2015)
	dec := EntireMonth(2015, 12)

	assert.True(t, year2015.Contains(dec))
	assert.True(t, dec.DoesNotContain(year2015))
}

func TestRange_Days(t *testing.T) {
	assert.Equal(t, 0, Empty().Days())
	assert.Equal(t, 1, OnlyToday().Days())
	assert.Equal(t, 365, EntireYear(2015).Days())
	assert.Equal(t, 366, EntireYear(2016).Days())
	assert.Equal(t, 29, EntireMonth(2016, 2).Days())
}

func TestRange_Error(t *testing.T) {
	assert.Nil(t, Never().Error())

	// Unbounded ranges are allowed
	end := Range{End: New(2015, 3, 1)}
	assert.Nil(t, end.Error(), "Unbounded start dates should not error")

	start := Range{Start: New(2015, 3, 1)}
	assert.Nil(t, start.Error(), "Unbounded end dates should not error")

	var invalid Range
	invalid.Start = New(2015, 3, 2)
	invalid.End = New(2015, 3, 1)
	assert.NotNil(t, invalid.Error())
}

func TestRange_Intersection(t *testing.T) {
	year2015 := EntireYear(2015)
	nov := EntireMonth(2015, 11)
	nov1 := New(2015, 11, 1)
	nov30 := New(2015, 11, 30)
	var zero Range

	assert.True(t, zero.Intersection(nov).IsZero())

	intersection := year2015.Intersection(nov)
	assert.Equal(t, nov1, intersection.Start)
	assert.Equal(t, nov30, intersection.End)
	assert.True(t, NewRange(nov1, nov30).Equals(intersection))
}

func TestRange_Marshal(t *testing.T) {
	// Empty ranges should render as null
	b, err := json.Marshal(Never())
	assert.Nil(t, err)
	assert.Equal(t, "null", string(b))

	// Infinite ranges should render as null start and end dates
	b, err = json.Marshal(Infinity())
	assert.Nil(t, err)
	assert.Equal(t, `{"start":null,"end":null}`, string(b))
}

func TestRange_Union(t *testing.T) {
	year2015 := EntireYear(2015)
	jan := EntireMonth(2016, 1)
	union := year2015.Union(jan)

	assert.Equal(t, New(2015, 1, 1), union.Start)
	assert.Equal(t, New(2016, 1, 31), union.End)

	feb := EntireMonth(2016, 2)
	union = jan.Union(feb)
	assert.Equal(t, feb.End, union.End)
	assert.Equal(t, jan.Start, union.Start)
}

func TestRange_Unmarshal(t *testing.T) {
	raw := `{"start":"2015-03-01","end":null}`
	var open Range
	assert.Nil(t, json.Unmarshal([]byte(raw), &open))
	assert.Equal(t, New(2015, 3, 1), open.Start)
	assert.True(t, open.End.IsZero())

	// TODO nulls should be unmarshaled as empty ranges
	// raw = `null`
	// var zero Range
	// assert.Nil(t, json.Unmarshal([]byte(raw), &zero))
	// assert.True(t, zero.IsEmpty())
}
