package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	macroPackageImportPath string
)

func init() {
	if len(os.Args) < 2 {
		fmt.Printf("error: a macro import path must be provided\n")
		fmt.Printf("usage: %s <import-path>\n", os.Args[0])
		os.Exit(1)
	}

	macroPackageImportPath = os.Args[1]
}

func main() {
	var (
		currentFilename         = os.Getenv("GOFILE")
		currentWorkingDirectory = os.Getenv("PWD")
		macroPackageName        = path.Base(macroPackageImportPath)
		fileSet                 = token.NewFileSet()
	)

	macroTagRegexp, err := regexp.Compile(fmt.Sprintf(`#\[%s\.(.*)\]`, macroPackageName))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
		return
	}

	astFile, err := parser.ParseFile(fileSet, currentFilename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("error: failed parsing %q: %s\n", currentFilename, err)
		os.Exit(4)
		return
	}

	for _, decl := range astFile.Decls {

		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {

			var (
				isGroupTypeSpec          = genDecl.Lparen.IsValid()
				macroFunctionInvocations = map[string]struct{}{}
			)

			if !isGroupTypeSpec {
				if genDecl.Doc != nil {
					for _, comment := range genDecl.Doc.List {
						if submatches := macroTagRegexp.FindStringSubmatch(comment.Text); len(submatches) > 1 {
							macroFunctionName := submatches[1]
							macroFunctionInvocations[macroFunctionName] = struct{}{}
						}
					}
				}
			}

			for _, spec := range genDecl.Specs {

				if typeSpec, ok := spec.(*ast.TypeSpec); ok {

					var (
						typeSpecName  = typeSpec.Name.Name
						startPosition = fileSet.Position(typeSpec.Pos())
						endPosition   = fileSet.Position(typeSpec.End())
					)

					if isGroupTypeSpec {
						if typeSpec.Doc != nil {
							for _, comment := range typeSpec.Doc.List {
								if submatches := macroTagRegexp.FindStringSubmatch(comment.Text); len(submatches) > 1 {
									macroFunctionName := submatches[1]
									macroFunctionInvocations[macroFunctionName] = struct{}{}
								}
							}
						}
					}

					// TODO: spawn a new craftMaterials:
					// 1. there would be more than one macro invocation
					// 2. to run in parallel

					for macroFunctionName, _ := range macroFunctionInvocations {

						macroOutputFilename := buildFilename(typeSpecName, macroPackageImportPath, macroFunctionName) + ".gen.go"

						programString, err := generateProgram(
							Values{
								Macro: Macro{
									Name:              macroFunctionName,
									PackageImportPath: macroPackageImportPath,
								},
								Source: Source{
									Filepath: filepath.Join(currentWorkingDirectory, currentFilename),
									StartPosition: Position{
										Line:   startPosition.Line,
										Column: startPosition.Column,
									},
									EndPosition: Position{
										Line:   endPosition.Line,
										Column: endPosition.Column,
									},
								},
								GenDecl: GenDecl{
									StartOffset: startPosition.Offset,
									EndOffset:   endPosition.Offset,
								},
								Output: Output{
									Filepath: filepath.Join(currentWorkingDirectory, macroOutputFilename),
								},
								Template: Template{
									TypeSpecName: typeSpecName,
									PackageName:  astFile.Name.Name,
								},
							},
						)
						if err != nil {
							fmt.Printf("error: failed generating program: %s\n", err)
							os.Exit(5)
							return
						}

						if err := craftMaterials(
							convertPathToFilename(currentWorkingDirectory),
							macroOutputFilename,
							programString,
						); err != nil {
							fmt.Println(err.Error())
							os.Exit(6)
							return
						}

					}

				}

			}

		}

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

	cmdOutput := &strings.Builder{}

	if gomodFile, err := os.Open(filepath.Join(tempPath, "go.mod")); errors.Is(err, os.ErrNotExist) {
		initCmd := exec.Command("go", "mod", "init", "macro")
		initCmd.Stdout = cmdOutput
		initCmd.Stderr = cmdOutput
		initCmd.Dir = tempPath
		if err := initCmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize macro module: %s", cmdOutput.String())
		}
	} else if err != nil {
		return fmt.Errorf("failed to open go.mod: %w", err)
	} else {
		if err := gomodFile.Close(); err != nil {
			return fmt.Errorf("failed to close go.mod: %w", err)
		}
	}

	cmdOutput.Reset()

	getCmd := exec.Command("go", "get", "-u")
	getCmd.Stdout = cmdOutput
	getCmd.Stderr = cmdOutput
	getCmd.Dir = tempPath
	if err := getCmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch macro dependencies: %s", cmdOutput.String())
	}

	cmdOutput.Reset()

	runMacroCmd := exec.Command("go", "run", filename)
	runMacroCmd.Stdout = cmdOutput
	runMacroCmd.Stderr = cmdOutput
	runMacroCmd.Dir = tempPath
	if err := runMacroCmd.Run(); err != nil {
		return fmt.Errorf("failed to run macro: %s", cmdOutput.String())
	}

	return nil
}

func buildFilename(structName string, macroPackageImportPath string, macroFunctionName string) string {
	return strings.ToLower(fmt.Sprintf("%s.%s.%s", structName, strings.ReplaceAll(macroPackageImportPath, "/", "."), macroFunctionName))
}

func convertPathToFilename(path string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(path, fmt.Sprintf("%c", filepath.Separator), "."),
		":",
		".",
	)
}
