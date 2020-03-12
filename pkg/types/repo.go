package types

import (
	"fmt"
	"strings"
)

// Repo struct
type Repo struct {
	Owner string
	Repo  string
}

// Repos for list
type Repos []Repo

// String
func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Repo)
}

// ComposeHTTPS compose https git address
func (r Repo) ComposeHTTPS() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", r.Owner, r.Repo)
}

// ComposeHTTPSWithCredential compose https git address with token
func (r Repo) ComposeHTTPSWithCredential(user, token string) string {
	return fmt.Sprintf("https://%s:%s@github.com/%s/%s.git",
		user, token, r.Owner, r.Repo)
}

// String
func (r Repos) String() string {
	var repoStrs []string
	for _, repo := range r {
		repoStrs = append(repoStrs, repo.String())
	}
	return strings.Join(repoStrs, ", ")
}
