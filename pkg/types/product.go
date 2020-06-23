package types

type ProductItem struct {
	Title    string
	Repo     Repo
	Children []ProductItem
}

// Product struct
type Product struct {
	Name      string
	Repos     []Repo
	Renames   map[Repo]Repo
	Structure []ProductItem
}
