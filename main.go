package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var NameDir = "testdata"

type file struct {
	Name    string
	size    int
	isDir   bool
	lastDir bool
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	filesDir := make([]file, 0)
	var dirs []string
	err := filepath.Walk(NameDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(out, "prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "skip" {
			return filepath.SkipDir
		}
		if path != NameDir {
			filesDir = append(filesDir, file{path, int(info.Size()), info.IsDir(), false})
			if info.IsDir() {
				dirs = append(dirs, info.Name())
			}

		}
		return nil

	})
	dirs = append(dirs, strings.Split(filesDir[0].Name, "/")[0])

	var maxLevel int
	for i := len(filesDir) - 1; i > 0; i-- {
		if filesDir[i].Name == "" {
			continue
		}
		parent := strings.Split(filesDir[i].Name, "/")[len(strings.Split(filesDir[i].Name, "/"))-2]
		if len(strings.Split(filesDir[i].Name, "/")) > maxLevel {
			maxLevel = strings.Count(filesDir[i].Name, "/")
		}
		index, contain := find(dirs, parent)
		if contain {
			dirs[index] = ""
			filesDir[i].lastDir = true
		}

	}

	levelPrint := make([]bool, maxLevel)
	for i := range levelPrint {
		levelPrint[i] = true
	}
	var str string
	for i, value := range filesDir {
		dir := strings.Split(value.Name, "/")
		if i > 0 && len(dir) < len(strings.Split(filesDir[i-1].Name, "/")) {
			for i := range levelPrint {
				levelPrint[i] = true
			}
		}
		levelIndex := 0
		for i := len(dir); i > 2; i-- {
			if levelPrint[levelIndex] {
				str += "│"
			}
			levelIndex++
			str += "\t"
		}
		if value.lastDir {
			str += "└───"
		} else {
			str += "├───"
		}
		str += dir[len(dir)-1]
		if !value.isDir {
			size := strconv.Itoa(value.size)
			if size == "0" {
				size = "empty"
			} else {
				size += "b"
			}
			str += " (" + size + ")"
		}
		str += "\n"
		if value.lastDir {
			levelPrint[strings.Count(value.Name, "/")-1] = false
		}
	}
	out.Write([]byte(str))
	return err
}

func find(arr []string, val string) (int, bool) {
	for i, value := range arr {
		if val == value {
			return i, true
		}
	}
	return -1, false
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
