package lib

import (
	"testing"
	"time"
)

func TestTimeElapsed(t *testing.T) {
	layout := "2006-01-02T15:04:05Z0700"

	currentDate := time.Date(2022, 01, 03,
		00, 00, 00, 00, time.UTC)

	tables := []struct {
		givenTime     string
		expectedYear  int
		expectedMonth int
		expectedDay   int
	}{
		{"2018-07-01T00:00:00Z", 3, 6, 2},
		// {"2016-02-29T00:00:00Z", 5, 10, 5},
		// {"2021-01-17T00:00:00Z", 2, 4, 3},
		// {"2021-01-17T00:00:00Z", 2, 7, 3},
	}

	for _, table := range tables {
		commitDate, err := time.Parse(layout, table.givenTime)
		if err != nil {
			t.Errorf("parse error occurred: %t", err)
		}

		year, month, day := TimeElapsed(currentDate, commitDate)

		if year != table.expectedYear || month != table.expectedMonth ||
			day != table.expectedDay {
			t.Errorf("outcome was incorrect,\n year - got: %d, want: %d\n month - got: %d, want %d\n day - got %d, want %d\n",
				year, table.expectedYear, month, table.expectedMonth, day, table.expectedDay)
		}
	}
}
