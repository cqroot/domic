package configs

import (
	"runtime"

	"github.com/cqroot/gmdots/pkg/file"
	"github.com/cqroot/gmdots/pkg/symlink"
)

func Apply(src, dest string) error {
	if runtime.GOOS == "windows" {
		return file.Apply(src, dest)
	}
	return symlink.Apply(src, dest)
}
