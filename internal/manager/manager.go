package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cqroot/domic/internal/utils"
	"gopkg.in/yaml.v3"
)

type DotfileConfig struct {
	Files map[string]string `json:"files"`
}

type Manager struct {
	workingDir string
	configFile string
	dotfiles   map[string]DotfileConfig
}

func New(configFile string) *Manager {
	return &Manager{
		workingDir: filepath.Dir(configFile),
		configFile: configFile,
		dotfiles:   make(map[string]DotfileConfig),
	}
}

func (m *Manager) LoadConfig() error {
	content, err := os.ReadFile(m.configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &m.dotfiles)
	if err != nil {
		return err
	}

	return nil
}

type CheckResult error

var (
	CheckResultOk             CheckResult = errors.New("OK")
	CheckResultSrcNotFound    CheckResult = errors.New("source file not found")
	CheckResultDstNotFound    CheckResult = errors.New("destination file not found")
	CheckResultFilesDifferent CheckResult = errors.New("files are different")
	CheckResultGetFileHashErr CheckResult = errors.New("get file hash error")
)

func (m *Manager) checkDotfile(name string, config DotfileConfig) CheckResult {
	for src, dst := range config.Files {
		src = filepath.Join(m.workingDir, name, src)
		dst = utils.ExpandPath(dst)

		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("%w: %s", CheckResultSrcNotFound, src)
		}

		if _, err := os.Stat(dst); err != nil {
			return fmt.Errorf("%w: %s", CheckResultDstNotFound, dst)
		}

		srcHash, err := utils.GetFileHash(src)
		if err != nil {
			return fmt.Errorf("%w: %s", CheckResultGetFileHashErr, err)
		}

		dstHash, err := utils.GetFileHash(dst)
		if err != nil {
			return fmt.Errorf("%w: %s", CheckResultGetFileHashErr, err)
		}

		if srcHash != dstHash {
			return fmt.Errorf("%w: %s -> %s", CheckResultFilesDifferent, src, dst)
		}
	}

	return CheckResultOk
}

func (m *Manager) Check() (map[string]CheckResult, error) {
	if err := m.LoadConfig(); err != nil {
		return nil, err
	}

	result := make(map[string]CheckResult)
	for name, config := range m.dotfiles {
		result[name] = m.checkDotfile(name, config)
	}

	return result, nil
}

func (m *Manager) Apply() error {
	if err := m.LoadConfig(); err != nil {
		return err
	}

	for name, config := range m.dotfiles {
		err := m.checkDotfile(name, config)
		if errors.Is(err, CheckResultOk) {
			continue
		}

		for src, dst := range config.Files {
			src := filepath.Join(m.workingDir, name, src)
			dst := utils.ExpandPath(dst)

			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return fmt.Errorf("error creating directory for %s: %v", name, err)
			}

			fmt.Printf("Applying %s -> %s\n", src, dst)
			if err := utils.CopyFile(src, dst); err != nil {
				return fmt.Errorf("error applying %s: %v", name, err)
			}
		}
	}
	return nil
}
