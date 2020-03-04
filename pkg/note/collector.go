package note

import (
	"github.com/google/go-github/v29/github"
)

// Collector for collect release notes
type Collector struct {
	github *github.Client
}

// New creates Collector instance
func New(g *github.Client) *Collector {
	return &Collector{g}
}
