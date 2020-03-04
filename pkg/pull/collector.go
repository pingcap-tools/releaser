package pull

import (
	"github.com/google/go-github/v29/github"
)

// Collector for collect pulls
type Collector struct {
	github *github.Client
}

// New creates Collector instance
func New(g *github.Client) *Collector {
	return &Collector{g}
}
