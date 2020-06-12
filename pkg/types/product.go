package types

// Product struct
type Product struct {
	Name    string
	Repos   []Repo
	Renames map[Repo]Repo
}
