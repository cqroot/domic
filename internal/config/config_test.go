package config_test

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cqroot/dotm/pkg/common"
	"github.com/cqroot/dotm/internal/config"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigRead(t *testing.T) {
	cfg, err := config.New("./testdata", "./testdata/dotm.toml")
	checkErr(t, err)

	baseDir, err := filepath.Abs("./testdata")
	checkErr(t, err)

	configDir, err := common.DotDir("config")
	checkErr(t, err)
	homeDir, err := common.DotDir("home")
	checkErr(t, err)

	expectedDots := []config.Dot{
		{
			Name:         "tmux",
			Source:       path.Join(baseDir, "tmux"),
			RelativePath: "tmux",
			Target:       path.Join(configDir, "tmux"),
			TargetType:   "config",
			Exec:         "tmux",
		}, {
			Name:         "vimrc",
			Source:       path.Join(baseDir, "vim/vimrc"),
			RelativePath: "vim/vimrc",
			Target:       path.Join(homeDir, ".vimrc"),
			TargetType:   "home",
			Exec:         "vim",
		},
	}

	assert.Equal(t, expectedDots, cfg.Dots)
}
