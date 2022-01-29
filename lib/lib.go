package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/iljaSL/dormant/model"
	"github.com/pterm/pterm"
)

func ReadFile(arg string) ([]model.Deps, error) {
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
		if strings.Contains(scanner.Text(), "git") {
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

// Getting the last recent commit page of repo with curl
//  curl \
//  -H "Accept: application/vnd.github.v3+json" \
//  https://api.github.com/repos/pterm/pterm/commits\?per_page\=1
// Time Format: 2021-11-17T14:20:35Z

// Getting the amount of commits since 6 moths with curl
//  curl \
//  -H "Accept: application/vnd.github.v3+json" \
//  https://api.github.com/repos/pterm/pterm/commits?since=2021-05-17T14:20:35Z
// Current Time - Config Time (default: 6 months) for 'since' time

// GetAPILastActivityInfo get the last recent project activity
// infos for go.mod dependencies
func GetAPILastActivityInfo(deps []model.Deps) ([]model.InspectResult, error) {
	result := []model.InspectResult{}

	for _, v := range deps {
		switch {
		case strings.Contains(v.URL, "github"):
			res, err := handleGitHubURL(v.Username, v.RepositoryName)
			if err != nil {
				return nil, err
			}

			result = append(result, model.InspectResult{
				URL:        v.URL,
				LastCommit: res[0].Commit.Author.Date,
			})
		case strings.Contains(v.URL, "gitlab"):
			// TODO HANDLE DATA WITH GITLAB API
			// https://stackoverflow.com/questions/39559689/where-do-i-find-the-project-id-for-the-gitlab-api
			// console.log(document.body.attributes[6])
			// https://gitlab.com/api/v4/projects/24747998
			// https://docs.gitlab.com/ee/api/projects.html#get-single-project
			//
		default:
			pterm.Error.Println("URL could not be handled")
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

		year, month, day := TimeElapsed(currentDate, commitDate)

		lastActivityDiff = append(lastActivityDiff, model.LastActivityDiff{
			URL:   v.URL,
			Year:  year,
			Month: month,
			Day:   day,
		})
	}

	return lastActivityDiff, nil
}

func TimeElapsed(a, b time.Time) (int, int, int) {
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

	return year, month, day
}
