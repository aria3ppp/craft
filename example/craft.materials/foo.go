package main

// import "%[1]s"
import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	macro "github.com/aria3ppp/craft/example/macros/hello"
)

func main() {
	// startPos, endPos := %[2]d, %[3]d
	startPos, endPos := int64(258), int64(327)

	// TODO: use the correct relative path
	// f, err := os.Open("%[4]s")
	f, err := os.Open("/home/alpha/workspace/go/craft/cmd/ast/ast.go")
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

	buffer := make([]byte, endPos-startPos)

	bytesRead, err := f.Read(buffer)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
		return
	}

	f.Close()

	structDef := string(buffer[:bytesRead])

	src := "package main\n\n" + structDef

	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, "snippet.go", src, parser.DeclarationErrors)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
		return
	}

	var genDecl *ast.GenDecl

	// TODO: is this correct?
	for _, decl := range file.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			genDecl = gd
		}
	}

	if genDecl == nil {
		if err != nil {
			fmt.Println("no *ast.StructType found")
			os.Exit(5)
			return
		}
	}

	// tmp, err := %[5]s.%[6]s(genDecl)
	templatedSourceCode, err := macro.Hello(genDecl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(6)
		return
	}

	tmp, err := template.New("program_template").Parse(templatedSourceCode)
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
			"PackageName":   "{{.PackageName}}",
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
