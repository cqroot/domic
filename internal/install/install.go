package install

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
	var errs []string

	for _, app := range apps {
		rawTarget, ok := config.TargetForOS(app)
		if !ok {
			fmt.Printf("[%s] skipped (no target for current OS)\n", app.Name)
			continue
		}
		target, err := config.ExpandHome(rawTarget)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: expand target: %v", app.Name, err))
			continue
		}
		source := filepath.Join(base, app.Path)
		if err := installEntry(app.Name, source, target); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("install completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func installEntry(name, source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source path %s does not exist", source)
		}
		return fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return walkAndInstall(name, source, target)
	}
	return processFile(name, source, target)
}

func walkAndInstall(name, sourceDir, targetDir string) error {
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
		return processFile(name, path, filepath.Join(targetDir, rel))
	})
}

func processFile(name, source, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	targetInfo, statErr := os.Stat(target)
	if statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat target: %w", statErr)
		}
		if err := copyFile(source, target); err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		fmt.Printf("[%s] installed %s\n", name, target)
		return nil
	}

	if targetInfo.IsDir() {
		return fmt.Errorf("target %s is a directory, expected file", target)
	}

	same, err := fileutil.SameContents(source, target)
	if err != nil {
		return err
	}
	if same {
		fmt.Printf("[%s] up to date %s\n", name, target)
		return nil
	}
	fmt.Printf("[%s] target already exists and differs from source\n  source: %s\n  target: %s\n", name, source, target)
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
