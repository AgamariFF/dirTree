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
	Name  string
	size  int
	isDir bool
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var check bool
	filesDir := make([]file, 2)
	err := filepath.Walk(NameDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(out, "prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "skip" {
			return filepath.SkipDir
		}
		if path != NameDir {
			filesDir = append(filesDir, file{path, int(info.Size()), info.IsDir()})
		}
		return nil

	})
	for i, file := range filesDir {
		filesDir[i].Name, _ = strings.CutPrefix(file.Name, NameDir)
	}
	for index, value := range filesDir {
		var prefStr,sufStr string
		sufStr = value.Name
		check = false
		if value.Name == "" {
			continue
		}
		count := strings.Count(value.Name, "\\")
		for count > 0 {
			if count > 1 {
				str, _, _ := strings.Cut(sufStr, "\\")
				prefStr += str + "\\"
				_, sufStr, _ = strings.Cut(sufStr, "\\")
				fmt.Println("prefStr: ", prefStr, ";sufStr: ", sufStr, "; Count:", count)
				 Если в будущем есть prefStr тогда "|"
			}
			if count == 1 {
				for i := index + 1; i < len(filesDir); i++ {
					if value.Name[:strings.LastIndex(value.Name, "\\")] == filesDir[i].Name[:strings.LastIndex(filesDir[i].Name, "\\")] {
						check = true
						fmt.Fprint(out, "├───")
						break
					}
				}
				if !check {
					fmt.Fprint(out, "└───")
				}
			}
			count--
		}
		fmt.Fprint(out, value.Name[strings.LastIndex(value.Name, "\\")+1:])
		if !value.isDir {
			if value.size > 0 {
				fmt.Fprint(out, " ("+strconv.Itoa(value.size)+"b)\n")
			} else {
				fmt.Fprint(out, " (empty)\n")
			}
		} else {
			fmt.Fprintln(out)
		}
		// fmt.Println(Name)
	}
	fmt.Println(filesDir)
	return err
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
