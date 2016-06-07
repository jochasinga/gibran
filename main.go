package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
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

func parseFile(path string, f os.FileInfo, err error) error {
	if strings.Contains(f.Name(), "broker") || strings.Contains(f.Name(), "@") {
		return nil
	}
	if !f.IsDir() && strings.Contains(f.Name(), ".go") {
		fset := token.NewFileSet()
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		conf := types.Config{Importer: importer.Default()}
		if err != nil {
			return err
		}
		base := filepath.Join(os.Getenv("GOPATH"), "src")
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		pkg, err := conf.Check(rel, fset, []*ast.File{file}, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Package: %q\n", pkg.Path())
		fmt.Printf("Name:    %s\n", pkg.Name())
		fmt.Printf("Imports: %s\n", pkg.Imports())
		fmt.Printf("Scope:   %s\n", pkg.Scope())
	}
	return nil
}

// readPackage is a WalkFunc for readProject
func readPackage(path string, f os.FileInfo, err error) error {

	// We want to ignore broker and delegate files and write to them
	// after having read other package files.
	if strings.Contains(f.Name(), "broker") || strings.Contains(f.Name(), "@") {
		return nil
	}

	if !f.IsDir() && strings.Contains(f.Name(), ".go") {
		// TODO: Parse file here and return an *Package
		// object that can be used to render template.
		fmt.Println(f.Name())
	}
	return nil
}

// readPackage reads sub packages of the project and create a
//relevant delegate and broker for each package.
func readProject(projectdir string) error {
	err := filepath.Walk(projectdir, readPackage)
	if err != nil {
		return err
	}
	err = filepath.Walk(projectdir, parseFile)
	if err != nil {
		return err
	}
	return nil
}

// createDir create a project directory structure based on the given projectMap
// TODO: Some schemes that might be interesting including, bare minimum, MVC, Flux, etc.
// It'd be great to provide users with a couple of structure options.
func createDir(projectName, projectDir string, projectMap map[string][]string) error {
	log.Printf("Creating %s...\n", projectName)
	for dir, files := range projectMap {
		fullpath := filepath.Join(projectDir, dir)
		log.Printf("Creating %s\n", fullpath)
		if strings.Contains(dir, ".go") {
			f, err := os.Create(fullpath)
			if err != nil {
				return err
			}
			defer f.Close()
			txt := []byte("package main")
			_, err = f.Write(txt)
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
			log.Printf("Creating %s\n", dst)
			if strings.Contains(file, ".go") {
				f, err := os.Create(dst)
				if err != nil {
					return err
				}
				defer f.Close()
				txt := []byte(fmt.Sprintf("package %s", dir))
				_, err = f.Write(txt)
				if err != nil {
					return err
				}
			} else {
				if err := os.Mkdir(dst, os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	log.Println("...Success!")
	return nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of Gibran CLI:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran startproject <projectName> <projectPath>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran run <projectPath>")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}
	flag.CommandLine.Init("", flag.ExitOnError)
}

func main() {
	flag.Parse()
	switch len(os.Args) {
	// No arguments provided
	case 1:
		fmt.Println("Supply a command...")
		flag.Usage()
		return
	// 1 arguments provided
	case 2:
		switch os.Args[1] {
		default:
			fmt.Println(os.Args[1])
			flag.Usage()
			return
		case "startproject":
			projectname := "myapp"
			rootdir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path := filepath.Join(rootdir, projectname)
			err = createDir(projectname, path, projectStructure)
			if err != nil {
				log.Fatal(err)
			}
			return
		case "run":
			rootdir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			if err := readProject(rootdir); err != nil {
				log.Fatal(err)
			}
			return
		}
	// 2 arguments provided
	case 3:
		switch os.Args[1] {
		default:
			fmt.Println(os.Args[1])
			flag.Usage()
			return
		case "startproject":
			projectname := os.Args[2]
			if projectname == "" {
				projectname = "myapp"
			}
			rootdir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path := filepath.Join(rootdir, projectname)
			err = createDir(projectname, path, projectStructure)
			if err != nil {
				log.Fatal(err)
			}
			return
		case "run":
			rootdir := os.Args[2]
			if rootdir == "" {
				root, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				rootdir = root
			}
			if err := readProject(rootdir); err != nil {
				log.Fatal(err)
			}
			return
		}

	// 3 arguments provided (likely with startproject <projectName> <projectPath>)
	case 4:
		if os.Args[1] != "startproject" {
			fmt.Println(os.Args[1])
			flag.Usage()
			return
		}
		projectname := os.Args[2]
		if projectname == "" {
			projectname = "myapp"
		}
		rootdir := os.Args[3]
		if rootdir == "" {
			root, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			rootdir = root
		}

		path := filepath.Join(rootdir, projectname)
		err := createDir(projectname, path, projectStructure)
		if err != nil {
			log.Fatal(err)
		}
		return
	default:
		flag.Usage()
		return
	}
}
