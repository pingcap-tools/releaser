package types

import "fmt"

// Repo struct
type Repo struct {
	Owner string
	Repo  string
}

// String
func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Repo)
}
