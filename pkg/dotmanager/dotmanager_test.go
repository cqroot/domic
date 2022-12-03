package dotmanager_test

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cqroot/dotm/pkg/common"
	"github.com/cqroot/dotm/pkg/dotmanager"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotsWithTag(t *testing.T) {
	baseDir, err := filepath.Abs("../../testdata")
	checkErr(t, err)

	configDir, err := common.DotDir("config")
	checkErr(t, err)

	dm, err := dotmanager.New("../../testdata", "../../testdata/dotm.toml")
	checkErr(t, err)

	dots := dm.DotsWithTag("term")

	expectedDots := []dotmanager.Dot{
		{
			Name:         "tmux",
			Source:       path.Join(baseDir, "tmux"),
			RelativePath: "tmux",
			Target:       path.Join(configDir, "tmux"),
			TargetType:   "config",
			Exec:         "tmux",
			Tags:         []string{"term"},
		},
	}

	assert.Equal(t, expectedDots, dots)
}
