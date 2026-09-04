package diffcmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/distribute"
	"github.com/cqroot/domic/internal/fileutil"
)

func Run(appNames []string) error {
	apps := config.Apps()
	if len(apps) == 0 {
		fmt.Println("no apps configured")
		return nil
	}

	nameSet := make(map[string]struct{}, len(appNames))
	for _, n := range appNames {
		nameSet[n] = struct{}{}
	}

	base := config.SourceBase()
	var errs []string
	matched := 0

	for _, app := range apps {
		if len(nameSet) > 0 {
			if _, ok := nameSet[app.Name]; !ok {
				continue
			}
		}
		matched++

		if !app.BinAvailable() {
			continue
		}
		rawTarget, ok := config.TargetForOS(app)
		if !ok {
			continue
		}
		target, err := config.ExpandHome(rawTarget)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: expand target: %v", app.Name, err))
			continue
		}
		source := filepath.Join(base, app.Path)
		if err := diffApp(app, source, target); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
		}
	}

	if len(nameSet) > 0 && matched == 0 {
		errs = append(errs, fmt.Sprintf("no apps matched: %v", appNames))
	}

	if len(errs) > 0 {
		return fmt.Errorf("diff completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func diffApp(app config.App, source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source path %s does not exist", source)
		}
		return fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return diffDir(app, source, target)
	}
	return diffFile(app, source, target)
}

func diffDir(app config.App, sourceDir, targetDir string) error {
	tmpDir, err := os.MkdirTemp("", "domic-diff-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

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
		target := filepath.Join(targetDir, rel)
		effectiveTarget := distribute.Target(app, path, target)
		sourcePath, ok, err := prepareDiff(app, path, effectiveTarget, tmpDir)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		return runDiff(sourcePath, effectiveTarget)
	})
}

func diffFile(app config.App, source, target string) error {
	tmpDir, err := os.MkdirTemp("", "domic-diff-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	sourcePath, ok, err := prepareDiff(app, source, target, tmpDir)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return runDiff(sourcePath, target)
}

// prepareDiff reports whether the target differs from the distributed bytes
// for `source`, and returns the path to pass as the source argument to the
// external diff command. For non-template sources the original source path is
// returned; for template sources the rendered bytes are written to a temp
// file and that path is returned.
func prepareDiff(app config.App, source, target, tmpDir string) (string, bool, error) {
	distributed, err := distribute.Build(app, source)
	if err != nil {
		return "", false, err
	}

	targetInfo, statErr := os.Stat(target)
	if statErr == nil {
		if targetInfo.IsDir() {
			return "", false, fmt.Errorf("target %s is a directory", target)
		}
		targetMD5, err := fileutil.MD5(target)
		if err != nil {
			return "", false, err
		}
		if distribute.MD5(distributed) == targetMD5 {
			return "", false, nil
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return "", false, fmt.Errorf("stat target: %w", statErr)
	}

	if distribute.IsTemplate(app, source) {
		tmpFile := filepath.Join(tmpDir, filepath.Base(source))
		if err := os.WriteFile(tmpFile, distributed, 0o644); err != nil {
			return "", false, fmt.Errorf("write temp: %w", err)
		}
		return tmpFile, true, nil
	}
	return source, true, nil
}

func runDiff(source, target string) error {
	cmdArgs := config.DiffCommand()
	if len(cmdArgs) == 0 {
		return fmt.Errorf("diff command not configured")
	}
	program, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		return fmt.Errorf("%s not found in PATH: %w", cmdArgs[0], err)
	}
	targetArg := target
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		targetArg = os.DevNull
	}
	args := append(append([]string{}, cmdArgs[1:]...), source, targetArg)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil
		}
		return err
	}
	return nil
}
