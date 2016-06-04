package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var projectStructure = map[string][]string{
	"brokers": []string{},
	"controllers": []string{"@.go", "basic.go"},
	"models": []string{"@.go", "user.go"},
	"routers": []string{"@.go", "routers.go"},
	"tests": []string{},
	"views": []string{"@.go", "views.go"},
	//"main.go": []string{},
}

func parseDir(projectDir string, projectMap map[string][]string) error {
	for dir, files := range projectMap {
		fullpath := filepath.Join(projectDir, dir)
		if dir == "main.go" {
			if _, err := os.Create(fullpath); err != nil {
				panic(err)
			}
		}
		if err := os.MkdirAll(fullpath, os.ModePerm); err != nil {
			return err
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

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of Gibran CLI:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran new <project_name>")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}
	
	flag.CommandLine.Init("", flag.ExitOnError)
}

func main() {
	commandname := flag.String("command", "", "Specify the command name.")
	projectname := flag.String("project", "", "Specify the project name.")
	rootdir := flag.String("root", "",  "Specify project's root directory.")

	flag.Parse()

	if *commandname == "" {
		flag.Usage()
		return
	}

	root := ""

	if *commandname == "startproject" {
		if *rootdir == "" {
			root, _ = os.Getwd()
			/*
			if err != nil {
				panic(err)
			}
                        */
		} else {
			root = *rootdir
		}
	}

	path := filepath.Join(root, *projectname)
	err := parseDir(path, projectStructure)
	if err != nil {
		panic(err)
	}
}
