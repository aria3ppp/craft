//go:generate craft github.com/aria3ppp/craft/example/macros/hello

package main

// This is a doc #1
// This is a doc #2
// This is a doc #3
import (
	"fmt"
	"log"
)

type X interface{}

// #[hello.Hello]
type (
	I interface{}
	T interface{}
)

// #[hello.Hello]
// #[hello.MarshalJSON]
// This is a doc #1
// This is a doc #2
// This is a doc #3
type PhonyStruct struct {
	PhoneyField string `json:"phoney_field"`
}

// #[hello.Hello]
// #[hello.MarshalJSON]
type XXX struct{}

type (
	// Abc struct
	// #[hello.Hello]
	// #[hello.MarshalJSON]
	Abc struct{}
	// #[hello.Hello]
	// #[hello.MarshalJSON]
	Xyz struct{}
)

func main() {
	// call generated Hello method
	var ps PhonyStruct
	fmt.Println(ps.Hello())
	jsonBytes, err := ps.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))
	fmt.Println()

	var xyz Xyz
	fmt.Println(xyz.Hello())
}
