/*
Craft bring macro meta programming to golang world!
To use a macro in your program add the following go:generate line and hash tags:

```go
//go:generate craft github.com/aria3ppp/craft/example/macros/hello

//#[hello.Hello]
type Name string

var john = Name("John")
john.Hello()
````

The 'hello' package (in this case 'github.com/aria3ppp/craft/example/macros/hello') is a special package defined somewhere on the internet
that have the following specification:
	- 'hello' package is behind 'macro' build tag
	- A function in hello package take in a '*ast.GenDecl' as input and return a `template.Template` plus an `error` as output

`craft` command build a program that will craft the source code by executing defined macros with passed `*ast.GenDecl` from source code as input.

Check this chat from chatgpt for more information about implementation design: https://chat.openai.com/c/8d42cd33-d6ba-4f91-8b15-30d41bd62130

You can use go template engine syntax extensively for many programmatical operations on the templated string: like 'for' loops or 'if' statements and other go template statements
*/

package craft
