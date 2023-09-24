package macro

import (
	"fmt"
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

func MarshalJSON(st *ast.StructType) (string, error) {
	var program string
	for _, field := range st.Fields.List {
		var name string
		for _, n := range field.Names {
			name += fmt.Sprintf(", %s", n.String())
		}
		name = name + ":"
		program += fmt.Sprintf("%s, %#v, %#v\n", name, field.Type, field.Tag)
	}
	return program, nil
}
