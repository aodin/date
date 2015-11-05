package date

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	day := New(2015, 3, 1)
	assert.Equal(t, "2015-03-01", day.String())

	output, err := json.Marshal(day)
	assert.Nil(t, err, "JSON marshaling of dates should not error")
	assert.Equal(t, []byte(`"2015-03-01"`), output)

	nextDay := New(2015, 3, 2)
	assert.True(t, nextDay.Equals(day.AddDays(1)))

	parsed, err := Parse("2015-03-01")
	assert.Nil(t, err, "Parsing of properly formatted dates should not error")
	assert.True(t, parsed.Equals(day))

	err = day.UnmarshalJSON([]byte(`"2015-03-01"`))
	assert.Nil(
		t, err, "UnmarshalJSON of a valid slice of bytes should not error",
	)

	jan1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.True(t, jan1.Before(day.Time))
}

func TestDate_Within(t *testing.T) {
	march1 := New(2015, 3, 1)

	feb := EntireMonth(2015, 2)
	march := EntireMonth(2015, 3)

	assert.False(t, march1.Within(feb))
	assert.Equal(t, march1, march.Start)
	assert.True(t, march1.Within(march))
}
