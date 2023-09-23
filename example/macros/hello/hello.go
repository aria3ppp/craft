package macro

import (
	"go/ast"
)

func Hello(st *ast.StructType) (string, error) {
	// structName := st.End()
	// st.Fields.List[0].Type

	// var x ast.TypeSpec
	// x.Type.

	programTemplate := `package {{.PackageName}}

	func (this *{{.StructName}}) Hello() string {
		return "Hello {{.StructName}}"
	}
	`

	return programTemplate, nil
}
