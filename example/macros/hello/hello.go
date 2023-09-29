package macro

import (
	"errors"
	"fmt"
	"go/ast"
)

func Hello(ts *ast.TypeSpec) (string, error) {
	_, isInterfaceType := ts.Type.(*ast.InterfaceType)

	if isInterfaceType {
		return "", errors.New("Hello do not work on interface types!")
	}

	programTemplate := `package {{.PackageName}}

func (this *{{.TypeSpecName}}) Hello() string {
	return "Hello {{.TypeSpecName}}"
}
`

	return programTemplate, nil
}

func MarshalJSON(ts *ast.TypeSpec) (string, error) {
	var structureString string

	switch typ := ts.Type.(type) {
	case *ast.ArrayType:
		structureString = "array[<donotknow>]"
	case *ast.ChanType:
		return "", errors.New("MarshalJSON do not work on channel types!")
	case *ast.FuncType:
		return "", errors.New("MarshalJSON do not work on function types!")
	case *ast.InterfaceType:
		return "", errors.New("MarshalJSON do not work on interface types!")
	case *ast.MapType:
		structureString = "map[<donknow>]<dontknow>"
	case *ast.StructType:
		for _, field := range typ.Fields.List {
			var name string
			for _, n := range field.Names {
				name += fmt.Sprintf(", %s", n.String())
			}
			name = name + ":"
			structureString += fmt.Sprintf("%s, %#v, %#v\n", name, field.Type, field.Tag)
		}
	}

	programTemplate := fmt.Sprintf(
		`package {{.PackageName}}
	
func (this *{{.TypeSpecName}}) MarshalJSON() ([]byte, error) {
	return []byte(%q), nil
}
`,
		structureString,
	)

	return programTemplate, nil
}
