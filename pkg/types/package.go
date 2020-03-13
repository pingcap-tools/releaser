package types

// Package ...
type Package struct {
	Name         string
	Repo         Repo
	Type         string
	URL          string
	Dependencies []Dependency
}

// Dependency ...
type Dependency struct {
	Name    string
	Version string
}
