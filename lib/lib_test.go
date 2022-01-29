package lib

import (
	"reflect"
	"testing"
	"time"

	"github.com/iljaSL/dormant/model"
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
		{"1969-09-15T00:00:00Z", 52, 3, 18},
		{"1692-09-15T00:00:00Z", 329, 3, 18},
		// Leap year is not considered yet in TimeElapsed
		// {"2016-02-29T00:00:00Z", 5, 10, 5},
	}

	tablesTwo := []struct {
		givenTime string
	}{
		{"2023-"},
		{"01T00:00:00Z"},
	}

	for _, table := range tables {
		commitDate, err := time.Parse(layout, table.givenTime)
		if err != nil {
			t.Errorf("parse error occurred: %t", err)
		}

		year, month, day, _ := TimeElapsed(currentDate, commitDate)
		if year != table.expectedYear || month != table.expectedMonth ||
			day != table.expectedDay {
			t.Errorf("outcome was incorrect,\n year - got: %d, want: %d\n month - got: %d, want %d\n day - got %d, want %d\n",
				year, table.expectedYear, month, table.expectedMonth, day, table.expectedDay)
		}
	}

	for _, table := range tablesTwo {
		commitDate, _ := time.Parse(layout, table.givenTime)

		_, _, _, err := TimeElapsed(currentDate, commitDate)
		if err == nil {
			t.Errorf("an error should have occurred")
		}
	}
}

func TestCalculateDepsActivity(t *testing.T) {
	type args struct {
		activityInfo []model.InspectResult
	}
	tests := []struct {
		name    string
		args    args
		want    []model.LastActivityDiff
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateDepsActivity(tt.args.activityInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDepsActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateDepsActivity() = %v, want %v", got, tt.want)
			}
		})
	}
}
