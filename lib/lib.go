package lib

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/iljaSL/dormant/model"
)

func ReadGoModFile(arg string) ([]model.Deps, error) {
	deps := []model.Deps{}

	f, err := os.Open(arg)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// read the file line by line
	scanner := bufio.NewScanner(f)

	reGitLink := regexp.MustCompile(`git[^\s]+`)
	reURLDetails := regexp.MustCompile(`[^\/\s]{1,}`)

	// Skip first line of go.mod, as it is the applications module name
	scanner.Scan()
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), model.GIT) {
			URLDetails := reURLDetails.FindAllString(scanner.Text(), 4)

			deps = append(deps, model.Deps{
				URL:            reGitLink.FindString(scanner.Text()),
				Username:       URLDetails[1],
				RepositoryName: URLDetails[2],
				Version:        URLDetails[3],
				Indirect:       strings.Contains(scanner.Text(), "indirect"),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return deps, err
}

// GetAPILastActivityInfo get the last recent project activity
// infos for go.mod dependencies
func GetAPILastActivityInfo(deps []model.Deps) ([]model.InspectResult, error) {
	result := []model.InspectResult{}

	for _, v := range deps {
		switch {
		case strings.Contains(v.URL, model.GITHUB):
			res, err := handleGitHubURL(v.Username, v.RepositoryName)
			if err != nil {
				return nil, err
			}

			result = append(result, model.InspectResult{
				URL:        v.URL,
				LastCommit: res[0].Commit.Author.Date,
			})
		default:
			return nil, fmt.Errorf("URL could not be handled")
		}
	}

	return result, nil
}

func handleGitHubURL(username, reponame string) ([]model.GitHubCommit, error) {
	gitHubCommitInfo := []model.GitHubCommit{}
	preparedURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=1",
		username, reponame)

	resp, err := http.Get(preparedURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("API rate limit exceeded for your IP address, authenticated requests feature is comming soon")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &gitHubCommitInfo)
	if err != nil {
		return nil, err
	}

	return gitHubCommitInfo, err
}

// CalculateDepsActivity calculate the dependency activity duration
func CalculateDepsActivity(activityInfo []model.InspectResult) ([]model.LastActivityDiff, error) {
	lastActivityDiff := []model.LastActivityDiff{}
	layout := "2006-01-02T15:04:05Z0700"
	t := time.Now()

	currentDate := time.Date(t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)

	for _, v := range activityInfo {
		commitDate, err := time.Parse(layout, v.LastCommit)
		if err != nil {
			return nil, err
		}

		year, month, err := TimeElapsed(currentDate, commitDate)
		if err != nil {
			return nil, err
		}

		month = month + (year * 12)

		lastActivityDiff = append(lastActivityDiff, model.LastActivityDiff{
			URL:   v.URL,
			Month: month,
		})
	}

	return lastActivityDiff, nil
}

func TimeElapsed(a, b time.Time) (int, int, error) {
	if b.IsZero() {
		return 0, 0, errors.New("wrong time format")
	}

	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year := int(y2 - y1)
	month := int(M2 - M1)
	day := int(d2 - d1)

	// Normalize negative values
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return year, month, nil
}
