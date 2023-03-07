package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cqroot/doter/pkg/dotfiles"
	"github.com/cqroot/doter/pkg/path"
)

func formattedDest(dest string) string {
	if dest == "" {
		return "-"
	}
	return "`" + strings.ReplaceAll(dest, "\\", "/") + "`"
}

func main() {
	path.ReplaceHomeDir("$HOME")
	path.ReplaceAppDataDir("%APPDATA%")
	path.ReplaceLocalAppDataDir("%LOCALAPPDATA%")

	dots := dotfiles.Dotfiles

	names := make([]string, 0, len(dots))
	for name := range dots {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("# Supported Applications")
	fmt.Println()

	for _, name := range names {
		dot := dots[name]

		fmt.Printf("## %s\n\n", name)

		fmt.Printf("- Linux: %s\n", formattedDest(dot.DstFunc("linux")))
		fmt.Printf("- MacOS: %s\n", formattedDest(dot.DstFunc("darwin")))
		fmt.Printf("- Windows: %s\n", formattedDest(dot.DstFunc("windows")))

		if dot.Doc != "" {
			fmt.Println("- Doc:", dot.Doc)
		}
		fmt.Println()
	}
}
