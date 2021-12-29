package model

// Deps go.mod dependency definition
type Deps struct {
	URL            string
	Username       string
	RepositoryName string
	Version        string
	Indirect       bool
}

// InspectResult result definition for the inspect flag
type InspectResult struct {
	URL        string
	LastCommit string
}

// GitHubCommit github commit API return definition
type GitHubCommit struct {
	Commit Commit `json:"commit"`
}

// Commit github sub commit definition
type Commit struct {
	Author Author `json:"author"`
}

// Author github author definition
type Author struct {
	Date string `json:"date"`
}
