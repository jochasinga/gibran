package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

type tmplConfig struct {
	Path       string
	Src        string
	Name       string
	ReadWriter io.ReadWriter
}

type tmplModel struct {
	Package *types.Package
	Info    *types.Info
}

// writeTmplWith takes a data struct and use the config information
// to write the corresponding template to the appropriate writer.
func writeTmplWith(data interface{}, config *tmplConfig) error {
	if config == nil {
		return errors.New("tmplConfig is nil")
	}
	t := template.New(config.Name)
	if config.Path != "" {
		tmp, err := t.ParseFiles(config.Path)
		if err != nil {
			return err
		}
		t = tmp
	} else {
		if config.Src == "" {
			return errors.New("tmplConfig: Path and Src empty")
		}
		tmp, err := t.Parse(config.Src)
		if err != nil {
			return err
		}
		t = tmp
	}
	if config.ReadWriter == nil {
		return errors.New("tmplConfig: Writer is nil")
	}
	err := t.Execute(config.ReadWriter, data)
	if err != nil {
		return err
	}
	return nil
}

func parseFile(path string, f os.FileInfo, err error) error {
	if strings.Contains(f.Name(), "broker") || strings.Contains(f.Name(), "@") {
		return nil
	}
	if !f.IsDir() && strings.Contains(f.Name(), ".go") {
		tmplModel := &tmplModel{}
		tmplConf := &tmplConfig{}
		conf := types.Config{Importer: importer.Default()}
		info := &types.Info{
			Defs: make(map[*ast.Ident]types.Object),
			Uses: make(map[*ast.Ident]types.Object),
		}

		fset := token.NewFileSet()
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		// Get the "github.com/jochasinga/packagename" portion
		prefix := filepath.Join(os.Getenv("GOPATH"), "src")
		rel, err := filepath.Rel(prefix, path)
		if err != nil {
			return err
		}
		pkg, err := conf.Check(rel, fset, []*ast.File{file}, info)
		if err != nil {
			// Just a warning print because even with an error
			// the pkg and info are not entirely nil.
			return err
		}

		// Test print attributes
		/*
			fmt.Printf("Package: %q\n", pkg.Path())
			fmt.Printf("Name:    %s\n", pkg.Name())
			fmt.Printf("Imports: %s\n", pkg.Imports())
			fmt.Printf("Scope:   %s\n", pkg.Scope())

			// Print out info
			for id, obj := range info.Defs {
				fmt.Printf("%s: %q DEFINES %v\n",
					fset.Position(id.Pos()), id.Name, obj)
			}
			for id, obj := range info.Uses {
				fmt.Printf("%s: %q USES %v\n",
					fset.Position(id.Pos()), id.Name, obj)
			}
		*/

		temp := `
                        {{with .Package}}
                        //!+ Path: "{{.Path}}"
                        package {{.Name}}

                        import (
                                {{range .Imports}}
                                  "{{.Name -}}"
                                {{end}}
                        )
                        {{end}}
                        {{with .Info}}
                        {{range $id, $obj := .Uses}}
                          {{$id}} {{$obj -}}
                        {{end}}
                        {{end}}
                        `
		tmplModel.Package = pkg
		tmplModel.Info = info
		tmplConf.Src = temp
		tmplConf.ReadWriter = new(bytes.Buffer)
		err = writeTmplWith(tmplModel, tmplConf)
		if err != nil {
			return err
		}
		fmt.Println(tmplConf.ReadWriter)
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

	// TODO: Skip the main package
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
