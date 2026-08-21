package apply

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cqroot/domic/internal/config"
)

type cache map[string]string

func cachePath() string {
	return filepath.Join(config.ConfigDir(), ".domic-cache.json")
}

func loadCache() cache {
	c := cache{}
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, &c)
	return c
}

func (c cache) save() error {
	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath(), data, 0o644)
}

func (c cache) get(target string) (string, bool) {
	v, ok := c[target]
	return v, ok
}

func (c cache) set(target, md5sum string) {
	c[target] = md5sum
}