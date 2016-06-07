package main

import (
	"bytes"
	"testing"
)

// mockPackage mocks types.Package with limited
// field functions  to simulate the method sets.
type mockPackage struct {
	Complete func() bool
	Imports  func() []*mockPackage
	Name     func() string
	Path     func() string
}

const temp = `
        //!+ Path: {{call .Path}}
        package {{call .Name}}

        import (
                {{range call .Imports}}
                {{call .Name -}}
                {{end}}
        )`

func TestWriteTemplate(t *testing.T) {
	expect := `
        //!+ Path: foo/bar/baz
        package melancholy

        import (

                packageA
                packageB
        )`
	config := &tmplConfig{
		Src:        temp,
		ReadWriter: new(bytes.Buffer),
	}
	data := &mockPackage{
		Complete: func() bool { return true },
		Name:     func() string { return "melancholy" },
		Path:     func() string { return "foo/bar/baz" },
		Imports: func() []*mockPackage {
			return []*mockPackage{
				&mockPackage{
					Name: func() string { return "packageA" },
					Path: func() string { return "a/b/c" },
				},
				&mockPackage{
					Name: func() string { return "packageB" },
					Path: func() string { return "d/e/f" },
				},
			}
		},
	}

	err := writeTmplWith(data, config)
	if err != nil {
		t.Error(err)
	}
	tmplStr := config.ReadWriter.(*bytes.Buffer).String()
	if tmplStr == expect+"\n" {
		t.Errorf("Unexpected: expect %s. got %s", expect, tmplStr)
	}
}
