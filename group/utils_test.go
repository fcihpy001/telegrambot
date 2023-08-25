package group

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertToSeconds(t *testing.T) {
	var ss []string = []string{
		"4",
		" 4 ",
		"4m",
		"4 m",
		"4min",
		" 4  min",
		"4minutes",

		"4h",
		"4hour",
		"4hours",

		"4d",
		"4day",
		"4Days",

		"4w",
		"4Week",
		"4WEEKS",

		"4month",
		"4Months",

		"4y",
		"4Year",
		"4Years",
	}
	var results = []int64{
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,

		4 * 3600,
		4 * 3600,
		4 * 3600,

		4 * 86400,
		4 * 86400,
		4 * 86400,

		4 * 86400 * 7,
		4 * 86400 * 7,
		4 * 86400 * 7,

		4 * 86400 * 30,
		4 * 86400 * 30,

		4 * 86400 * 365,
		4 * 86400 * 365,
		4 * 86400 * 365,
	}
	for i, s := range ss {
		seconds := convertToSeconds(s)
		assert.Equal(t, results[i], seconds)
	}
}

func TestParseTime(t *testing.T) {
	ss := []struct {
		str string
		tm  int64
	}{
		{"2023-08-25 14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{" 2023-08-25   14:30:00 ", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-8-25 14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-8-25 14:30", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		// {"14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-08-25", time.Date(2023, 8, 25, 0, 0, 0, 0, time.Local).Unix()},
	}

	for _, item := range ss {
		val, err := parseTime(item.str, true)
		assert.Nil(t, err)
		assert.Equal(t, item.tm, val)
	}
}
