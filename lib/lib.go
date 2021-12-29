package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

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

// Todo: Only URL is not enough, need to get Username and Repo name from the URL
// Todo: to be able to use it with the GitHub API

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

// Todo: Better Name for function, as it only handles the commit date get info
func GetAPIInfo(deps []model.Deps) error {
	result := []model.InspectResult{}

	for _, v := range deps {
		switch {
		case strings.Contains(v.URL, "github"):
			res, _ := handleGitHubURL(v.Username, v.RepositoryName)

			result = append(result, model.InspectResult{
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
	fmt.Println("TEMP RES LOG", result)
	return nil
}

func handleGitHubURL(username, reponame string) ([]model.GitHubCommit, error) {
	preparedURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=1",
		username, reponame)

	resp, err := http.Get(preparedURL)
	if err != nil {
		log.Fatalln(err)
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data []model.GitHubCommit
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
