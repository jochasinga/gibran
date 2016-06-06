package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var projectStructure = map[string][]string{
	"brokers":     []string{},
	"controllers": []string{"@.go", "basic.go"},
	"models":      []string{"@.go", "user.go"},
	"routers":     []string{"@.go", "routers.go"},
	"tests":       []string{},
	"views":       []string{"@.go", "views.go"},
	"main.go":     []string{},
}

func createDir(projectDir string, projectMap map[string][]string) error {
	for dir, files := range projectMap {
		fullpath := filepath.Join(projectDir, dir)
		if strings.Contains(dir, ".go") {
			_, err := os.Create(fullpath)
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(fullpath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		for _, file := range files {
			dst := filepath.Join(fullpath, file)
			if strings.Contains(file, ".go") {
				if _, err := os.Create(dst); err != nil {
					return err
				}
			} else {
				if err := os.Mkdir(dst, os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

var (
	commandname = flag.String("command", "", "Specify the command name.")
	projectname = flag.String("project", "", "Specify the project name.")
	rootdir     = flag.String("root", "", "Specify project's root directory.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of Gibran CLI:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran startproject <projectName> <projectPath>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran run")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}
	flag.CommandLine.Init("", flag.ExitOnError)
}

func main() {
	flag.Parse()
	if *commandname == "" {
		flag.Usage()
		return
	}
	if *commandname == "run" {
		// parse project's package
		// create and update brokers
		// run project with go run
	}
	if *commandname == "startproject" {
		if *rootdir == "" {
			root, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			*rootdir = root
		}
	}
	path := filepath.Join(*rootdir, *projectname)
	err := createDir(path, projectStructure)
	if err != nil {
		panic(err)
	}
}
