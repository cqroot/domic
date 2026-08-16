package status

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/fileutil"
)

type FileState int

const (
	FileOK FileState = iota
	FileMissing
	FileModified
)

func (s FileState) String() string {
	switch s {
	case FileOK:
		return "ok"
	case FileMissing:
		return "missing"
	case FileModified:
		return "modified"
	}
	return "unknown"
}

type AppState int

const (
	AppOK AppState = iota
	AppMissing
	AppModified
	AppPartial
	AppSkipped
)

func (s AppState) String() string {
	switch s {
	case AppOK:
		return "ok"
	case AppMissing:
		return "missing"
	case AppModified:
		return "modified"
	case AppPartial:
		return "partial"
	case AppSkipped:
		return "skipped"
	}
	return "unknown"
}

type FileResult struct {
	Source string
	Target string
	State  FileState
}

type AppResult struct {
	Name   string
	Source string
	Target string
	State  AppState
	Files  []FileResult
}

func Run() error {
	apps := config.Apps()
	if len(apps) == 0 {
		fmt.Println("no apps configured")
		return nil
	}

	base := config.SourceBase()
	var results []AppResult
	var errs []string

	for _, app := range apps {
		rawTarget, ok := config.TargetForOS(app)
		if !ok {
			results = append(results, AppResult{
				Name:   app.Name,
				Source: filepath.Join(base, app.Path),
				State:  AppSkipped,
			})
			continue
		}
		target, err := config.ExpandHome(rawTarget)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: expand target: %v", app.Name, err))
			continue
		}
		source := filepath.Join(base, app.Path)
		result, err := inspect(app.Name, source, target)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
			continue
		}
		results = append(results, result)
	}

	printResults(results)

	if len(errs) > 0 {
		return fmt.Errorf("status completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func inspect(name, source, target string) (AppResult, error) {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AppResult{}, fmt.Errorf("source path %s does not exist", source)
		}
		return AppResult{}, fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return inspectDir(name, source, target)
	}
	return inspectFile(name, source, target)
}

func inspectFile(name, source, target string) (AppResult, error) {
	state, err := compareFile(source, target)
	if err != nil {
		return AppResult{}, err
	}
	return AppResult{
		Name:   name,
		Source: source,
		Target: target,
		State:  appStateForSingle(state),
	}, nil
}

func inspectDir(name, sourceDir, targetDir string) (AppResult, error) {
	var files []FileResult
	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, walkErr error) error {
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
		state, err := compareFile(path, filepath.Join(targetDir, rel))
		if err != nil {
			return err
		}
		files = append(files, FileResult{
			Source: path,
			Target: filepath.Join(targetDir, rel),
			State:  state,
		})
		return nil
	})
	if err != nil {
		return AppResult{}, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Target < files[j].Target
	})
	return AppResult{
		Name:   name,
		Source: sourceDir,
		Target: targetDir,
		State:  aggregateFileStates(files),
		Files:  files,
	}, nil
}

func compareFile(source, target string) (FileState, error) {
	targetInfo, err := os.Stat(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return FileMissing, nil
		}
		return 0, fmt.Errorf("stat target: %w", err)
	}
	if targetInfo.IsDir() {
		return 0, fmt.Errorf("target %s is a directory", target)
	}
	same, err := fileutil.SameContents(source, target)
	if err != nil {
		return 0, err
	}
	if same {
		return FileOK, nil
	}
	return FileModified, nil
}

func appStateForSingle(s FileState) AppState {
	switch s {
	case FileOK:
		return AppOK
	case FileMissing:
		return AppMissing
	case FileModified:
		return AppModified
	}
	return 0
}

func aggregateFileStates(files []FileResult) AppState {
	if len(files) == 0 {
		return AppMissing
	}
	var ok, missing, modified int
	for _, f := range files {
		switch f.State {
		case FileOK:
			ok++
		case FileMissing:
			missing++
		case FileModified:
			modified++
		}
	}
	if modified > 0 {
		return AppModified
	}
	if ok == len(files) {
		return AppOK
	}
	if ok > 0 && missing > 0 {
		return AppPartial
	}
	return AppMissing
}

func printResults(results []AppResult) {
	fmt.Printf("%s %s %s\n",
		paintHeader(nameColWidth, "NAME"),
		paintHeader(statusColWidth, "STATUS"),
		paintHeader(0, "TARGET"),
	)
	for _, r := range results {
		fmt.Printf("%-*s %s %s\n", nameColWidth, r.Name, paintAppState(r.State), r.Target)
		for _, f := range r.Files {
			fmt.Printf("%-*s %s %s\n", nameColWidth, "  ", paintFileState(f.State), f.Target)
		}
	}
}

const (
	nameColWidth   = 15
	statusColWidth = 10

	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiGray   = "\033[90m"
)

var colorEnabled = os.Getenv("NO_COLOR") == ""

func paint(color, s string) string {
	if !colorEnabled {
		return s
	}
	return color + s + ansiReset
}

func paintHeader(width int, s string) string {
	return paint(ansiBold, fmt.Sprintf("%-*s", width, s))
}

func paintAppState(s AppState) string {
	return paint(appStateColor(s), fmt.Sprintf("%-*s", statusColWidth, s.String()))
}

func paintFileState(s FileState) string {
	return paint(fileStateColor(s), fmt.Sprintf("%-*s", statusColWidth, s.String()))
}

func appStateColor(s AppState) string {
	switch s {
	case AppOK:
		return ansiGreen
	case AppMissing, AppPartial:
		return ansiYellow
	case AppModified:
		return ansiRed
	case AppSkipped:
		return ansiGray
	}
	return ""
}

func fileStateColor(s FileState) string {
	switch s {
	case FileOK:
		return ansiGreen
	case FileMissing:
		return ansiYellow
	case FileModified:
		return ansiRed
	}
	return ""
}
