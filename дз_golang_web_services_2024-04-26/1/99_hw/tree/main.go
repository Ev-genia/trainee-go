package main

import (
	// "io"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	fork        = "├───"
	branch      = "│	"
	endOfBranch = "└───"
)

type FileInfo struct {
	BasePath  string
	CleanPath string
	IsDir     bool
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {
	files, dirs, err := getRows(path, printFiles)
	if err != nil {
		return err
	}
	if printFiles {
		err = printRows(files, out)
	} else {
		err = printRows(dirs, out)
	}
	return err
}

func getRows(path string, printFiles bool) ([]FileInfo, []FileInfo, error) {
	var files []FileInfo
	var dirs []FileInfo
	var err error
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if printFiles && !info.IsDir() {
			size := "empty"
			if info.Size() != 0 {
				size = strconv.FormatInt(info.Size(), 10)
			}
			files = append(files, FileInfo{BasePath: filepath.Base(path), CleanPath: filepath.Clean(path) + " (" + size + ")", IsDir: false})
		}

		if info.IsDir() {
			dirs = append(dirs, FileInfo{BasePath: filepath.Base(path), CleanPath: filepath.Clean(path), IsDir: true})
			files = append(files, FileInfo{BasePath: filepath.Base(path), CleanPath: filepath.Clean(path), IsDir: true})
		}
		return nil
	})

	return files, dirs, err
}

func printRows(rows []FileInfo, out *os.File) (err error) {
	slashesPreview := 1
	for i, row := range rows {
		slashes := strings.Count(row.CleanPath, "/")
		// fmt.Println(row.CleanPath)
		if slashes > 0 {
			if slashes >= slashesPreview {
				fmt.Fprintln(out, strings.Repeat(branch, slashes-1)+fork+row.BasePath)
			} else {
				if i < len(rows)-1 {
					fmt.Fprintln(out, fork+row.BasePath)
				} else {
					fmt.Fprintln(out, endOfBranch+row.BasePath)
				}
			}
		}
		slashesPreview = slashes
	}

	return err
}
