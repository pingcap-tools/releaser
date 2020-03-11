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

// String
func (r Repos) String() string {
	var repoStrs []string
	for _, repo := range r {
		repoStrs = append(repoStrs, repo.String())
	}
	return strings.Join(repoStrs, ", ")
}
