package apply

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/fileutil"
)

func Run() error {
	apps := config.Apps()
	if len(apps) == 0 {
		fmt.Println("no apps configured")
		return nil
	}

	base := config.SourceBase()
	cache := loadCache()
	var errs []string

	for _, app := range apps {
		rawTarget, ok := config.TargetForOS(app)
		if !ok {
			fmt.Printf("[%s] %s (no target for current OS)\n", app.Name, paint(ansiGray, "skipped"))
			continue
		}
		target, err := config.ExpandHome(rawTarget)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: expand target: %v", app.Name, err))
			continue
		}
		source := filepath.Join(base, app.Path)
		if err := applyEntry(app.Name, source, target, cache); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
		}
	}

	if err := cache.save(); err != nil {
		errs = append(errs, fmt.Sprintf("save cache: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("apply completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func applyEntry(name, source, target string, cache cache) error {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source path %s does not exist", source)
		}
		return fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return walkAndApply(name, source, target, cache)
	}
	return processFile(name, source, target, cache)
}

func walkAndApply(name, sourceDir, targetDir string, cache cache) error {
	return filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		return processFile(name, path, filepath.Join(targetDir, rel), cache)
	})
}

func processFile(name, source, target string, cache cache) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	sourceMD5, err := fileutil.MD5(source)
	if err != nil {
		return fmt.Errorf("hash source: %w", err)
	}

	targetInfo, statErr := os.Stat(target)
	if statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat target: %w", statErr)
		}
		if err := copyFile(source, target); err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		cache.set(target, sourceMD5)
		fmt.Printf("[%s] %s %s\n", name, paint(ansiGreen, "applied"), target)
		return nil
	}

	if targetInfo.IsDir() {
		return fmt.Errorf("target %s is a directory, expected file", target)
	}

	targetMD5, err := fileutil.MD5(target)
	if err != nil {
		return fmt.Errorf("hash target: %w", err)
	}

	if targetMD5 == sourceMD5 {
		cache.set(target, sourceMD5)
		fmt.Printf("[%s] %s %s\n", name, paint(ansiGreen, "up to date"), target)
		return nil
	}

	if recorded, ok := cache.get(target); ok && recorded == targetMD5 {
		if err := copyFile(source, target); err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		cache.set(target, sourceMD5)
		fmt.Printf("[%s] %s %s\n", name, paint(ansiGreen, "applied"), target)
		return nil
	}

	fmt.Printf("[%s] %s\n  source: %s\n  target: %s\n", name, paint(ansiYellow, "target already exists and differs from source"), source, target)
	return nil
}

func copyFile(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
	ansiGreen = "\033[32m"
	ansiYellow = "\033[33m"
	ansiGray  = "\033[90m"
)

var colorEnabled = os.Getenv("NO_COLOR") == ""

func paint(color, s string) string {
	if !colorEnabled {
		return s
	}
	return color + s + ansiReset
}
