package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	macroImportPath     string
	macroOutputFilename string
)

func init() {
	flag.StringVar(&macroImportPath, "by", "", "macro import path")
	flag.StringVar(&macroOutputFilename, "to", "", "macro output filename")
	flag.Parse()

	if macroImportPath == "" {
		fmt.Printf("error: invalid value for flag -by: %q\n", macroImportPath)
		os.Exit(1)
	}

	if macroOutputFilename == "" {
		fmt.Printf("error: invalid value for flag -to: %q\n", macroOutputFilename)
		os.Exit(1)
	}
}

func main() {
	printEnv()
	defer fmt.Println()

	goFile := os.Getenv("GOFILE")

	pkgName := path.Base(macroImportPath)

	hashTagRegexp, err := regexp.Compile(fmt.Sprintf(`#\[%s\.(.*)\]`, pkgName))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
		return
	}

	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, goFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("error: failed parsing %q file: %s\n", goFile, err)
		os.Exit(4)
		return
	}

	var (
		macroDefinitionName string
		programString       string
	)

	for _, decl := range file.Decls {

		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE {

			for _, comment := range gd.Doc.List {
				if submatches := hashTagRegexp.FindStringSubmatch(comment.Text); len(submatches) > 1 {
					macroDefinitionName = submatches[1]

				}
			}

			for _, spec := range gd.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok {

						var err error
						programString, err = generateProgram(
							Values{
								Macro: Macro{
									Name:               macroDefinitionName,
									CraftedProgramPath: filepath.Join(os.Getenv("PWD"), goFile),
									PackageImportPath:  macroImportPath,
								},
								GenDecl: GenDecl{
									StartOffset: fs.Position(gd.Pos()).Offset,
									EndOffset:   fs.Position(gd.End()).Offset,
								},
								OutputFilename: filepath.Join(os.Getenv("PWD"), macroOutputFilename),
								StructName:     typeSpec.Name.Name,
								PackageName:    file.Name.Name,
							},
						)
						if err != nil {
							fmt.Printf("error: failed generating program: %s\n", err)
							os.Exit(5)
							return
						}

					}
				}
			}

		}
	}

	// no hash tags matched macro
	if programString == "" {
		return
	}

	if err := craftMaterials(pkgName, macroDefinitionName+".go", programString); err != nil {
		fmt.Println(err.Error())
		os.Exit(6)
		return
	}

}

var craftTempDir = "craft"

func craftMaterials(subDir string, filename string, content string) error {
	tempPath := filepath.Join(os.TempDir(), craftTempDir, subDir)

	if err := os.MkdirAll(tempPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filepath.Join(tempPath, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to the file: %w", err)
	}

	stdErr := &strings.Builder{}

	if gomodFile, err := os.Open(filepath.Join(tempPath, "go.mod")); errors.Is(err, os.ErrNotExist) {
		initCmd := exec.Command("go", "mod", "init", "macro")
		initCmd.Stderr = stdErr
		initCmd.Dir = tempPath
		if err := initCmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize macro module: %w\nerror message:\n%s", err, stdErr.String())
		}
	} else if err != nil {
		return fmt.Errorf("failed to open go.mod: %w", err)
	} else {
		if err := gomodFile.Close(); err != nil {
			return fmt.Errorf("failed to close go.mod: %w", err)
		}
	}

	stdErr.Reset()

	getCmd := exec.Command("go", "get", "-u")
	getCmd.Stderr = stdErr
	getCmd.Dir = tempPath
	if err := getCmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch macro dependencies: %w\nerror message:\n%s", err, stdErr.String())
	}

	stdErr.Reset()

	runMacroCmd := exec.Command("go", "run", filename)
	runMacroCmd.Stderr = stdErr
	runMacroCmd.Dir = tempPath
	if err := runMacroCmd.Run(); err != nil {
		return fmt.Errorf("failed to run macro: %w\nerror message:\n%s", err, stdErr.String())
	}

	return nil
}

func printEnv() {
	_, b, _, _ := runtime.Caller(0)
	fmt.Println("_, b, _, _ := runtime.Caller(0) =", b)
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(`out, err := exec.Command("go", "env", "GOMOD").Output() =`, string(out))
	fmt.Println("\tPWD =", os.Getenv("PWD"))
	fmt.Println("\tGOARCH =", os.Getenv("GOARCH"))
	fmt.Println("\tGOOS =", os.Getenv("GOOS"))
	fmt.Println("\tGOFILE =", os.Getenv("GOFILE"))
	fmt.Println("\tGOLINE =", os.Getenv("GOLINE"))
	fmt.Println("\tGOPACKAGE =", os.Getenv("GOPACKAGE"))
	fmt.Println("\tGOROOT =", os.Getenv("GOROOT"))
	fmt.Println("\tDOLLAR =", os.Getenv("DOLLAR"))
	fmt.Println("\tPATH =", os.Getenv("PATH"))
}
