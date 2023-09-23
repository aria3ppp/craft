package main

import (
	"html/template"
	"strings"
)

var programTemplate = `package main

import (
	"strings"
	"text/template"

	macro "{{.Macro.PackageImportPath}}"
)

func main() {
	startPos, endPos := int64({{.GenDecl.StartOffset}}), int64({{.GenDecl.EndOffset}})

	// TODO: use the correct relative path
	f, err := os.Open("{{.Macro.CraftedProgramPath}}")
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

	src := "package main\n\n"+ structDef

	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, "snippet.go", src, parser.DeclarationErrors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
		return
	}

	var structType *ast.StructType

	// TODO: is this correct?
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if st, isStructType := spec.(*ast.TypeSpec).Type.(*ast.StructType); isStructType {
					structType = st
				}
			}
		}
	}

	programTemplate, err := macro.{{.Macro.Name}}(structType)
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
			"DeclTokenName": "{{.DeclTokenName}}",
			"PackageName": "{{.PackageName}}",
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

	// TODO: make '*.gen.go' file off the template
}

func craft(programString string) error {
	fmt.Println(programString)
	return nil
}
`

type Values struct {
	Macro   Macro
	GenDecl GenDecl

	DeclTokenName string
	PackageName   string
}

type Macro struct {
	Name               string
	CraftedProgramPath string
	PackageImportPath  string
}

type GenDecl struct {
	StartOffset int
	EndOffset   int
}

func generateProgram(values Values, outputPath string) (string, error) {
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
