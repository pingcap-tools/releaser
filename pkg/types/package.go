package types

// Package ...
type Package struct {
	Name         string
	Repo         Repo
	Dependencies []Dependency
}

// Dependency ...
type Dependency struct {
	Name    string
	Version string
}
