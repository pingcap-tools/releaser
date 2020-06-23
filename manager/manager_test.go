package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/you06/releaser/pkg/types"
)

func TestParseStructure(t *testing.T) {
	s, e := parseStructure([]string{"pingcap/tidb", "tikv/tikv", "pingcap/pd", "Tools: pingcap/br, pingcap/dumpling, pingcap/tidb-lightning, pingcap/ticdc"})
	assert.Nil(t, e)
	assert.Equal(t, s, []types.ProductItem{
		{
			Repo: types.Repo{
				Owner: "pingcap",
				Repo:  "tidb",
			},
		},
		{
			Repo: types.Repo{
				Owner: "tikv",
				Repo:  "tikv",
			},
		},
		{
			Repo: types.Repo{
				Owner: "pingcap",
				Repo:  "pd",
			},
		},
		{
			Title: "Tools",
			Children: []types.ProductItem{
				{
					Repo: types.Repo{
						Owner: "pingcap",
						Repo:  "br",
					},
				},
				{
					Repo: types.Repo{
						Owner: "pingcap",
						Repo:  "dumpling",
					},
				},
				{
					Repo: types.Repo{
						Owner: "pingcap",
						Repo:  "tidb-lightning",
					},
				},
				{
					Repo: types.Repo{
						Owner: "pingcap",
						Repo:  "ticdc",
					},
				},
			},
		},
	})
}
