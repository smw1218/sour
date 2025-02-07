package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func (sd *ServiceData) modifyProxyList(filename string) error {
	// Read and parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %v", err)
	}

	// Add new import
	importPath := fmt.Sprintf("github.com/smw1218/sour/cmd/%s-service/app", sd.ServiceName)
	addImport(node, sd.PackageAlias(), importPath)

	// Modify AllServices slice
	sd.modifyAllServices(node)

	// Create custom configuration for formatting
	printConfig := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 8,
	}

	// Write the modified AST back to file with custom formatting
	var buf bytes.Buffer
	if err := printConfig.Fprint(&buf, fset, node); err != nil {
		return fmt.Errorf("error formatting node: %v", err)
	}

	// Post-process the output to ensure slice elements are on separate lines
	output := formatSliceLiteral(buf.String())
	//output := buf.String()

	// gofmt
	formatted, err := format.Source([]byte(output))
	if err != nil {
		return fmt.Errorf("error gofmt: %w", err)
	}
	return os.WriteFile(filename, formatted, 0644)
}

func formatSliceLiteral(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	inSlice := false
	indent := ""

	for _, line := range lines {
		if strings.Contains(line, "AllServices = []") {
			inSlice = true
			indent = strings.Repeat("\t", strings.Count(line, "\t")+1)

			openBrace := strings.Index(line, "{")
			result = append(result, line[:openBrace+1])
			// assume this is the first insertion
			if len(line[openBrace+1:]) > 0 {
				result = append(result, indent+line[openBrace+1:len(line)-1]+",", "}")
			}
			continue
		}

		if inSlice && strings.Contains(line, "}") {
			inSlice = false
			elements := strings.Split(strings.Trim(strings.TrimSuffix(line, "}"), " \t"), ",")
			for _, element := range elements {
				element = strings.TrimSpace(element)
				if element != "" {
					result = append(result, indent+element+",")
				}
			}
			result = append(result, strings.TrimRight(line[:strings.Index(line, "}")], " \t")+"}")
			continue
		}
		if inSlice {
			result = append(result, singleLineLiterals(line, indent)...)
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func singleLineLiterals(line string, indent string) []string {
	elems := strings.Split(line, ",")
	if len(elems) < 2 {
		return []string{line}
	}
	lines := make([]string, 0, len(elems))
	for _, e := range elems {
		v := strings.TrimSpace(e)
		if v != "" {
			lines = append(lines, indent+v+",")
		}
	}
	return lines
}

func addImport(node *ast.File, serviceName, importPath string) {
	// Find the import declaration
	var importDecl *ast.GenDecl
	for _, decl := range node.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			importDecl = gen
			break
		}
	}

	// Create new import spec
	newImport := &ast.ImportSpec{
		Name: ast.NewIdent(serviceName),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s"`, importPath),
		},
	}

	if importDecl != nil {
		importDecl.Specs = append(importDecl.Specs, newImport)
	} else {
		// Create new import declaration if none exists
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{newImport},
		}
		node.Decls = append([]ast.Decl{importDecl}, node.Decls...)
	}
}

func (sd *ServiceData) modifyAllServices(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == "AllServices" {
						// Create the new service call expression
						serviceCall := &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(sd.PackageAlias()),
								Sel: ast.NewIdent(fmt.Sprintf("New%sService", sd.TitleName())),
							},
						}

						// If there's no existing value, create a new composite literal
						if len(valueSpec.Values) == 0 {
							valueSpec.Values = []ast.Expr{
								&ast.CompositeLit{
									Type: valueSpec.Type,
									Elts: []ast.Expr{serviceCall},
								},
							}
						} else {
							// Add to existing composite literal
							if comp, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
								comp.Elts = append(comp.Elts, serviceCall)
							}
						}
					}
				}
			}
		}
		return true
	})
}
