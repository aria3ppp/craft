//go:generate craft github.com/aria3ppp/craft/example/macros/hello

package main

// This is a doc #1
// This is a doc #2
// This is a doc #3
import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

// #[hello.Hello]
// This is a doc #1
// This is a doc #2
// This is a doc #3
type PhonyStruct struct {
	PhoneyField string `json:"phoney_field"`
}

// #[hello.Hello]
type XXX struct{}

// #[hello.Hello]
type (
	// Abc struct
	// #[hello.Hello]
	Abc struct{}
	// #[hello.Hello]
	Xyz struct{}
)

func main() {
	// call generated Hello method
	var ps PhonyStruct
	fmt.Println(ps.Hello())
	fmt.Println()

	var xyz Xyz
	xyz.Hello()
	fmt.Println()

	// Create a new file set.
	fs := token.NewFileSet()

	// Parse the Go source code file.
	// Replace "your_program.go" with the path to your Go source file.
	file, err := parser.ParseFile(fs, "cmd/ast/ast.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Print the AST of the parsed Go program.
	printAST(file, fs)
}

func printAST(node interface{}, fs *token.FileSet) {
	// This function recursively prints the AST nodes.
	// You can customize the format as needed.
	switch n := node.(type) {
	case *ast.File:
		fmt.Println("File:", fs.Position(n.Pos()), fs.Position(n.End()))
		for _, decl := range n.Decls {
			printAST(decl, fs)
		}
	case *ast.GenDecl:
		fmt.Println("Generic Declaration:", fs.Position(n.Pos()), fs.Position(n.End()))
		fmt.Println("Doc:", n.Doc.List[0].Text)
		fmt.Println("Token:", n.Tok)
		for _, spec := range n.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				fmt.Println("Type:", fs.Position(n.Pos()), fs.Position(n.End()))
				fmt.Println("Name:", ts.Name.Name)
				printAST(ts.Type, fs)
			}
		}
	case *ast.StructType:
		fmt.Println("Struct:", fs.Position(n.Pos()), fs.Position(n.End()))
		for _, field := range n.Fields.List {
			fmt.Println("Field:", fs.Position(n.Pos()), fs.Position(n.End()))
			fmt.Println("Name:", field.Names)
			fmt.Println("FieldType:", field.Type)
			// TODO: get more information off the field
		}
		// TODO: get more information off the struct type
	case *ast.FuncDecl:
		fmt.Println("Function:", fs.Position(n.Pos()), fs.Position(n.End()))
		fmt.Println("Name:", n.Name.Name)
		printAST(n.Body, fs)
		// You can print more details about the function here.
	case *ast.BlockStmt:
		fmt.Println("Block:", fs.Position(n.Pos()), fs.Position(n.End()))
		for _, stmt := range n.List {
			printAST(stmt, fs)
		}
	// Add more cases for other AST node types as needed.
	default:
		fmt.Printf("%T\n", n)
	}

	fmt.Println()
}

func _() {
	{
	}
}
