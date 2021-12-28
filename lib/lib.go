package lib

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
)

type Deps struct {
	URL            string
	Username       string
	RepositoryName string
	Version        string
	Indirect       bool
}

func ReadFile(arg string) []Deps {
	deps := []Deps{}

	f, err := os.Open(arg)
	if err != nil {
		pterm.Error.Println(err)
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

			deps = append(deps, Deps{
				URL:            reGitLink.FindString(scanner.Text()),
				Username:       URLDetails[1],
				RepositoryName: URLDetails[2],
				Version:        URLDetails[3],
				Indirect:       strings.Contains(scanner.Text(), "indirect"),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		pterm.Error.Println(err)
	}

	return deps
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
