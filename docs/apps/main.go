package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cqroot/doter/pkg/dotmanager"
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
	linuxDotMap := dotmanager.DefaultDotMap("linux")
	darwinDotMap := dotmanager.DefaultDotMap("darwin")
	windowsDotMap := dotmanager.DefaultDotMap("windows")

	names := make([]string, 0, len(linuxDotMap))
	for name := range linuxDotMap {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("# Supported Applications")
	fmt.Println()

	for _, name := range names {
		fmt.Printf("## %s\n\n", name)

		fmt.Printf("- Linux: %s\n", formattedDest(linuxDotMap[name].Dest))
		fmt.Printf("- MacOS: %s\n", formattedDest(darwinDotMap[name].Dest))
		fmt.Printf("- Windows: %s\n", formattedDest(windowsDotMap[name].Dest))

		if linuxDotMap[name].Doc != "" {
			fmt.Println("- Doc:", linuxDotMap[name].Doc)
		}
		fmt.Println()
	}
}
