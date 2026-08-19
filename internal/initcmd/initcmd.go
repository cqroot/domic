package initcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

func Run(arg string) error {
	target := filepath.Join(xdg.ConfigHome, "domic")
	repo := resolveRepo(arg)
	return clone(repo, target)
}

func resolveRepo(arg string) string {
	if isFullURL(arg) {
		return arg
	}
	if !strings.Contains(arg, "/") {
		return "https://github.com/" + arg + "/dotfiles"
	}
	return "https://github.com/" + arg
}

func isFullURL(s string) bool {
	return strings.Contains(s, "://") || strings.HasPrefix(s, "git@")
}

func clone(repo, target string) error {
	cmd := exec.Command("git", "clone", repo, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}
