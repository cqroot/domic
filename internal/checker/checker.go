/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package checker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cqroot/domic/internal/config"
	"github.com/cqroot/domic/internal/utils"
)

type FileCheckError error

var (
	FileCheckOk             FileCheckError = errors.New("OK")
	FileCheckSrcNotFound    FileCheckError = errors.New("source file not found")
	FileCheckDstNotFound    FileCheckError = errors.New("destination file not found")
	FileCheckFilesDifferent FileCheckError = errors.New("files are different")
	FileCheckGetFileHashErr FileCheckError = errors.New("get file hash error")
)

type FileCheckResult struct {
	Src   string
	Dst   string
	Error FileCheckError
}

func CheckFile(src, dst string) FileCheckResult {
	fileCheckResult := FileCheckResult{
		Src: src,
		Dst: dst,
	}
	if _, err := os.Stat(src); err != nil {
		fileCheckResult.Error = fmt.Errorf("%w: %s", FileCheckSrcNotFound, src)
		return fileCheckResult
	}

	if _, err := os.Stat(dst); err != nil {
		fileCheckResult.Error = fmt.Errorf("%w: %s", FileCheckDstNotFound, dst)
		return fileCheckResult
	}

	srcHash, err := utils.GetFileHash(src)
	if err != nil {
		fileCheckResult.Error = fmt.Errorf("%w: %s", FileCheckGetFileHashErr, err)
		return fileCheckResult
	}

	dstHash, err := utils.GetFileHash(dst)
	if err != nil {
		fileCheckResult.Error = fmt.Errorf("%w: %s", FileCheckGetFileHashErr, err)
		return fileCheckResult
	}

	if srcHash != dstHash {
		fileCheckResult.Error = fmt.Errorf("%w: %s -> %s", FileCheckFilesDifferent, src, dst)
		return fileCheckResult
	}

	fileCheckResult.Error = FileCheckOk
	return fileCheckResult
}

type DotfileResult struct {
	Name        string
	OkCnt       int
	TotalCnt    int
	FileResults []FileCheckResult
}

func CheckDotfile(workDir, name string, dotfile config.Dotfile) DotfileResult {
	dotfileResult := DotfileResult{
		Name:        name,
		OkCnt:       0,
		TotalCnt:    0,
		FileResults: make([]FileCheckResult, 0),
	}
	ch := make(chan FileCheckResult)

	wgCheck := sync.WaitGroup{}
	for src, dst := range dotfile.Files {
		wgCheck.Add(1)
		go func() {
			defer wgCheck.Done()

			ch <- CheckFile(filepath.Join(workDir, name, src), dst)
		}()
	}

	wgResult := sync.WaitGroup{}
	wgResult.Add(1)
	go func() {
		defer wgResult.Done()

		for result := range ch {
			dotfileResult.TotalCnt++
			if errors.Is(result.Error, FileCheckOk) {
				dotfileResult.OkCnt++
			} else {
				dotfileResult.FileResults = append(dotfileResult.FileResults, result)
			}
		}
	}()

	wgCheck.Wait()
	close(ch)
	wgResult.Wait()
	return dotfileResult
}

func CheckDotfiles(workDir string, dotfiles map[string]config.Dotfile) []DotfileResult {
	dotfilesResult := make([]DotfileResult, 0)
	ch := make(chan DotfileResult)

	wgCheck := sync.WaitGroup{}
	for name, dotfile := range dotfiles {
		wgCheck.Add(1)
		go func(name string) {
			defer wgCheck.Done()

			ch <- CheckDotfile(workDir, name, dotfile)
		}(name)
	}

	wgResult := sync.WaitGroup{}
	wgResult.Add(1)
	go func() {
		defer wgResult.Done()

		for result := range ch {
			dotfilesResult = append(dotfilesResult, result)
		}
	}()

	wgCheck.Wait()
	close(ch)
	wgResult.Wait()
	return dotfilesResult
}
