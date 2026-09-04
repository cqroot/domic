package status

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/distribute"
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

type Options struct {
	ShowFiles bool
}

type AppResult struct {
	Name   string
	Source string
	Target string
	State  AppState
	Files  []FileResult
}

func Run(opts Options) error {
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
		result, err := inspect(app, source, target)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", app.Name, err))
			continue
		}
		if !opts.ShowFiles {
			result.Files = nil
		}
		results = append(results, result)
	}

	printResults(results)

	if len(errs) > 0 {
		return fmt.Errorf("status completed with errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func inspect(app config.App, source, target string) (AppResult, error) {
	info, err := os.Stat(source)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AppResult{}, fmt.Errorf("source path %s does not exist", source)
		}
		return AppResult{}, fmt.Errorf("stat source %s: %w", source, err)
	}
	if info.IsDir() {
		return inspectDir(app, source, target)
	}
	return inspectFile(app, source, target)
}

func inspectFile(app config.App, source, target string) (AppResult, error) {
	effectiveTarget := distribute.Target(app, source, target)
	state, err := compareDistributed(app, source, effectiveTarget)
	if err != nil {
		return AppResult{}, err
	}
	return AppResult{
		Name:   app.Name,
		Source: source,
		Target: effectiveTarget,
		State:  appStateForSingle(state),
	}, nil
}

func inspectDir(app config.App, sourceDir, targetDir string) (AppResult, error) {
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
		target := filepath.Join(targetDir, rel)
		effectiveTarget := distribute.Target(app, path, target)
		state, err := compareDistributed(app, path, effectiveTarget)
		if err != nil {
			return err
		}
		files = append(files, FileResult{
			Source: path,
			Target: effectiveTarget,
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
		Name:   app.Name,
		Source: sourceDir,
		Target: targetDir,
		State:  aggregateFileStates(files),
		Files:  files,
	}, nil
}

func compareDistributed(app config.App, source, target string) (FileState, error) {
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
	distributed, err := distribute.Build(app, source)
	if err != nil {
		return 0, err
	}
	targetMD5, err := fileutil.MD5(target)
	if err != nil {
		return 0, err
	}
	if distribute.MD5(distributed) == targetMD5 {
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
		fmt.Printf("%-*s %s %s\n", nameColWidth, r.Name, paintAppState(r.State), paintTarget(r.Source, r.Target))
		for _, f := range r.Files {
			fmt.Printf("%-*s %s %s\n", nameColWidth, "  ", paintFileState(f.State), paintTarget(f.Source, f.Target))
		}
	}
}

func paintTarget(source, target string) string {
	prefix := commonPathPrefix(source, target)
	if prefix == "" {
		return target
	}
	return paint(ansiCyan, prefix) + strings.TrimPrefix(target, prefix)
}

func commonPathPrefix(a, b string) string {
	aParts := strings.Split(filepath.Clean(a), string(filepath.Separator))
	bParts := strings.Split(filepath.Clean(b), string(filepath.Separator))
	n := len(aParts)
	if len(bParts) < n {
		n = len(bParts)
	}
	i := 0
	for i < n && aParts[i] == bParts[i] {
		i++
	}
	if i == 0 {
		return ""
	}
	return strings.Join(aParts[:i], string(filepath.Separator))
}

const (
	nameColWidth   = 15
	statusColWidth = 10

	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
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
