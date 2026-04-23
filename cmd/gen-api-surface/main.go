// Command gen-api-surface emits a machine-readable snapshot of the public API
// of github.com/hra42/openrouter-go to docs/api-surface.json.
//
// The snapshot is intended for AI coding agents and doc tooling that want to
// answer "does this SDK have X?" with a single file read instead of grepping.
//
// Run from the repo root:
//
//	go run ./cmd/gen-api-surface
//
// CI should run this and fail if the result differs from the committed file.
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type apiSurface struct {
	Module   string     `json:"module"`
	Package  string     `json:"package"`
	Overview string     `json:"overview,omitempty"`
	Consts   []constVar `json:"consts,omitempty"`
	Vars     []constVar `json:"vars,omitempty"`
	Types    []typeDecl `json:"types,omitempty"`
	Funcs    []funcDecl `json:"funcs,omitempty"`
}

type constVar struct {
	Name string `json:"name"`
	Doc  string `json:"doc,omitempty"`
}

type typeDecl struct {
	Name    string     `json:"name"`
	Doc     string     `json:"doc,omitempty"`
	Kind    string     `json:"kind"` // struct | interface | alias | basic
	Methods []funcDecl `json:"methods,omitempty"`
}

type funcDecl struct {
	Name      string `json:"name"`
	Receiver  string `json:"receiver,omitempty"`
	Signature string `json:"signature"`
	Doc       string `json:"doc,omitempty"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen-api-surface:", err)
		os.Exit(1)
	}
}

func run() error {
	fset := token.NewFileSet()
	matches, err := filepath.Glob("*.go")
	if err != nil {
		return fmt.Errorf("glob: %w", err)
	}

	files := make([]*ast.File, 0, len(matches))
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if f.Name.Name != "openrouter" {
			// Ignore files that happen to live here but belong to a
			// different package (e.g. build-tag-scoped alternatives).
			continue
		}
		files = append(files, f)
	}
	if len(files) == 0 {
		return fmt.Errorf("no openrouter package files found in current directory (run from repo root)")
	}

	pkg, err := doc.NewFromFiles(fset, files, "github.com/hra42/openrouter-go", doc.PreserveAST)
	if err != nil {
		return fmt.Errorf("doc: %w", err)
	}

	out := apiSurface{
		Module:   "github.com/hra42/openrouter-go",
		Package:  pkg.Name,
		Overview: firstSentence(pkg, pkg.Doc),
	}

	for _, c := range pkg.Consts {
		for _, name := range c.Names {
			out.Consts = append(out.Consts, constVar{Name: name, Doc: firstSentence(pkg, c.Doc)})
		}
	}
	for _, v := range pkg.Vars {
		for _, name := range v.Names {
			out.Vars = append(out.Vars, constVar{Name: name, Doc: firstSentence(pkg, v.Doc)})
		}
	}

	for _, t := range pkg.Types {
		td := typeDecl{
			Name: t.Name,
			Doc:  firstSentence(pkg, t.Doc),
			Kind: typeKind(t),
		}
		for _, m := range t.Methods {
			td.Methods = append(td.Methods, funcDecl{
				Name:      m.Name,
				Receiver:  m.Recv,
				Signature: formatSignature(fset, m.Decl),
				Doc:       firstSentence(pkg, m.Doc),
			})
		}
		// Constructors appear under t.Funcs in go/doc.
		for _, f := range t.Funcs {
			out.Funcs = append(out.Funcs, funcDecl{
				Name:      f.Name,
				Signature: formatSignature(fset, f.Decl),
				Doc:       firstSentence(pkg, f.Doc),
			})
		}
		sort.Slice(td.Methods, func(i, j int) bool { return td.Methods[i].Name < td.Methods[j].Name })
		out.Types = append(out.Types, td)
	}
	for _, f := range pkg.Funcs {
		out.Funcs = append(out.Funcs, funcDecl{
			Name:      f.Name,
			Signature: formatSignature(fset, f.Decl),
			Doc:       firstSentence(pkg, f.Doc),
		})
	}

	sort.Slice(out.Consts, func(i, j int) bool { return out.Consts[i].Name < out.Consts[j].Name })
	sort.Slice(out.Vars, func(i, j int) bool { return out.Vars[i].Name < out.Vars[j].Name })
	sort.Slice(out.Types, func(i, j int) bool { return out.Types[i].Name < out.Types[j].Name })
	sort.Slice(out.Funcs, func(i, j int) bool { return out.Funcs[i].Name < out.Funcs[j].Name })

	if err := os.MkdirAll("docs", 0o755); err != nil {
		return err
	}
	path := filepath.Join("docs", "api-surface.json")
	buf, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	buf = append(buf, '\n')
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s (%d types, %d funcs)\n", path, len(out.Types), len(out.Funcs))
	return nil
}

func typeKind(t *doc.Type) string {
	if t.Decl == nil || len(t.Decl.Specs) == 0 {
		return "unknown"
	}
	ts, ok := t.Decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return "unknown"
	}
	switch ts.Type.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.FuncType:
		return "func"
	case *ast.Ident, *ast.SelectorExpr, *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.StarExpr:
		return "alias"
	default:
		return "other"
	}
}

func formatSignature(fset *token.FileSet, decl *ast.FuncDecl) string {
	if decl == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("func ")
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		sb.WriteString("(")
		sb.WriteString(exprString(decl.Recv.List[0].Type))
		sb.WriteString(") ")
	}
	sb.WriteString(decl.Name.Name)
	sb.WriteString(paramsString(decl.Type.Params))
	if decl.Type.Results != nil && len(decl.Type.Results.List) > 0 {
		sb.WriteString(" ")
		sb.WriteString(paramsString(decl.Type.Results))
	}
	return sb.String()
}

func paramsString(fl *ast.FieldList) string {
	if fl == nil {
		return "()"
	}
	parts := make([]string, 0, len(fl.List))
	for _, f := range fl.List {
		typ := exprString(f.Type)
		if len(f.Names) == 0 {
			parts = append(parts, typ)
			continue
		}
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
		parts = append(parts, strings.Join(names, ", ")+" "+typ)
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

func exprString(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + exprString(x.X)
	case *ast.SelectorExpr:
		return exprString(x.X) + "." + x.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprString(x.Elt)
	case *ast.MapType:
		return "map[" + exprString(x.Key) + "]" + exprString(x.Value)
	case *ast.ChanType:
		return "chan " + exprString(x.Value)
	case *ast.Ellipsis:
		return "..." + exprString(x.Elt)
	case *ast.FuncType:
		return "func" + paramsString(x.Params)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.IndexExpr:
		return exprString(x.X) + "[" + exprString(x.Index) + "]"
	case *ast.IndexListExpr:
		parts := make([]string, 0, len(x.Indices))
		for _, ix := range x.Indices {
			parts = append(parts, exprString(ix))
		}
		return exprString(x.X) + "[" + strings.Join(parts, ", ") + "]"
	default:
		return fmt.Sprintf("%T", e)
	}
}

func firstSentence(pkg *doc.Package, s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// go/doc synopsis: first sentence, up to newline or period+space.
	return strings.TrimSpace(pkg.Synopsis(s))
}
