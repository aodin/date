package date

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateRange(t *testing.T) {
	assert.True(t, Infinity().IsInfinity())
	assert.True(t, Empty().IsEmpty())

	assert.True(t, Never().Equals(Empty()))
	assert.True(t, Forever().Equals(Infinity()))
	assert.False(t, Forever().Equals(Empty()))
}

func TestDateRange_Contains(t *testing.T) {
	year2015 := EntireYear(2015)
	dec := EntireMonth(2015, 12)

	assert.True(t, year2015.Contains(dec))
	assert.True(t, dec.DoesNotContain(year2015))
}

func TestDateRange_Error(t *testing.T) {
	assert.Nil(t, Never().Error())

	var invalid DateRange
	invalid.Start = New(2015, 3, 2)
	invalid.End = New(2015, 3, 1)
	assert.NotNil(t, invalid.Error())
}

func TestDateRange_Intersection(t *testing.T) {
	year2015 := EntireYear(2015)
	nov := EntireMonth(2015, 11)
	var zero DateRange

	assert.True(t, zero.Intersection(nov).IsZero())

	intersection := year2015.Intersection(nov)
	assert.Equal(t, New(2015, 11, 1), intersection.Start)
	assert.Equal(t, New(2015, 11, 30), intersection.End)
}

func TestDateRange_Union(t *testing.T) {
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
