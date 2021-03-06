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
	}{
		{"2018-07-01T00:00:00Z", 3, 6},
		{"1969-09-15T00:00:00Z", 52, 3},
		{"1692-09-15T00:00:00Z", 329, 3},
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

		year, month, _ := TimeElapsed(currentDate, commitDate)
		if year != table.expectedYear || month != table.expectedMonth {
			t.Errorf("outcome was incorrect,\n year - got: %d, want: %d\n month - got: %d want: %d\n",
				year, table.expectedYear, month, table.expectedMonth)
		}
	}

	for _, table := range tablesTwo {
		commitDate, _ := time.Parse(layout, table.givenTime)

		_, _, err := TimeElapsed(currentDate, commitDate)
		if err == nil {
			t.Errorf("an error should have occurred")
		}
	}
}

func TestReadGoModFile(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Deps
		wantErr bool
	}{
		{
			name: "correct go mod file read",
			args: args{
				"./test.mod",
			},
			want: []model.Deps{
				{
					URL:            "github.com/spf13/viper",
					Username:       "spf13",
					RepositoryName: "viper",
					Version:        "v1.11.0",
					Indirect:       false,
				},
			},
		},
		{
			name:    "file does not exist",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadGoModFile(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAPILastActivityInfo(t *testing.T) {
	type args struct {
		deps []model.Deps
	}
	tests := []struct {
		name    string
		args    args
		want    []model.InspectResult
		wantErr bool
	}{
		{
			name:    "empty struct error",
			want:    []model.InspectResult{},
			wantErr: true,
		},
		{
			name: "check non supported url",
			args: args{
				deps: []model.Deps{
					{
						URL: "doesnotexist.com",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "check faulty github repo",
			args: args{
				deps: []model.Deps{
					{
						URL: "github.com/sdsd/viadsaper",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetAPILastActivityInfo(tt.args.deps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAPILastActivityInfo() = %v,\n want %v\n", got, tt.want)
			}
		})
	}
}

func Test_handleGitHubURL(t *testing.T) {
	type args struct {
		username string
		reponame string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.GitHubCommit
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				username: "test",
				reponame: "testRepo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleGitHubURL(tt.args.username, tt.args.reponame)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleGitHubURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleGitHubURL() = %v, want %v", got, tt.want)
			}
		})
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
		{
			name: "empty struct",
			args: args{},
			want: []model.LastActivityDiff{},
		},
		{
			name: "wring time format",
			args: args{
				[]model.InspectResult{
					{
						URL:        "github.com/test/test",
						LastCommit: time.Now().String(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CalculateDepsActivity(tt.args.activityInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDepsActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
