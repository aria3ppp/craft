package main

import (
	"html/template"
	"strings"
)

var programTemplate = `package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	macro "{{.Macro.PackageImportPath}}"
)

func main() {
	startPos, endPos := int64({{.GenDecl.StartOffset}}), int64({{.GenDecl.EndOffset}})

	f, err := os.Open("{{.Source.Filepath}}")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	if _, err = f.Seek(startPos, 0); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
		return
	}

	buffer := make([]byte, endPos - startPos)

	bytesRead, err := f.Read(buffer)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
		return
	}

	f.Close()

	structDef := string(buffer[:bytesRead])

	var src string

	if {{.NeedTypePrepend}} {
		src = "package main\n\ntype "+ structDef
	} else {
		src = "package main\n\n"+ structDef
	}

	fs := token.NewFileSet()

	astFile, err := parser.ParseFile(fs, "", src, parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
		return
	}

	var typeSpec *ast.TypeSpec

	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if ts, isTypeSpec := spec.(*ast.TypeSpec); isTypeSpec {
					typeSpec = ts
				}
			}
		}
	}

	programTemplate, err := macro.{{.Macro.Name}}(typeSpec)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
		return
	}

	tmp, err := template.New("program_template").Parse(programTemplate)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(6)
		return
	}

	stringBuilder := &strings.Builder{}

	if err := tmp.Execute(
		stringBuilder,
		map[string]any{
			"TypeSpecName": "{{.Template.TypeSpecName}}",
			"PackageName": "{{.Template.PackageName}}",
		},
	); err != nil {
		fmt.Println(err.Error())
		os.Exit(7)
		return
	}

	if err := craft(stringBuilder.String()); err != nil {
		fmt.Println(err.Error())
		os.Exit(8)
		return
	}
}

func craft(programString string) error {
	file, err := os.Create("{{.Output.Filepath}}")
	if err != nil {
		return err
	}

	if _, err := file.WriteString(programString); err != nil {
		return err
	}
	
	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
`

type Values struct {
	Macro           Macro
	Source          Source
	GenDecl         GenDecl
	Output          Output
	Template        Template
	NeedTypePrepend bool
}

type Macro struct {
	Name              string
	PackageImportPath string
}

type Source struct {
	Filepath string
}

type GenDecl struct {
	StartOffset int
	EndOffset   int
}

type Output struct {
	Filepath string
}

type Template struct {
	TypeSpecName string
	PackageName  string
}

func generateProgram(values Values) (string, error) {
	tmp, err := template.New("program").Parse(programTemplate)
	if err != nil {
		return "", err
	}

	var stringBuilder strings.Builder

	if err = tmp.Execute(&stringBuilder, values); err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}
