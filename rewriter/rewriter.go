package rewriter

import (
	"go/ast"
	"strings"

	"github.com/fatih/astrewrite"
)

type Rewriter interface {
	Rewrite(file *ast.File) *ast.File
}

type identRewriter struct {
	rewrites map[*ast.Object]ast.Expr
}

func NewIdentRewriter() Rewriter {
	return new(identRewriter)
}

func (r *identRewriter) rewriteIdents(n ast.Node) (ast.Node, bool) {
	if _, ok := n.(*ast.ValueSpec); ok {
		return n, false
	}

	v, ok := n.(*ast.Ident)
	if !ok || !isRewrite(v) {
		return n, true
	}

	return r.rewrites[v.Obj], false
}

func (r *identRewriter) rewriteDeclarations(n ast.Node) bool {
	v, ok := n.(*ast.ValueSpec)
	if !ok {
		return true
	}

	k := len(v.Values)
	if k < 1 {
		return false
	}

	for i := 0; i < k; i++ {
		v.Values[i] = astrewrite.Walk(v.Values[i], r.rewriteIdents).(ast.Expr)
		if name := v.Names[i]; isRewrite(name) {
			r.rewrites[name.Obj] = v.Values[i]
		}
	}

	return false
}

func (r *identRewriter) Rewrite(file *ast.File) *ast.File {
	r.rewrites = make(map[*ast.Object]ast.Expr)
	ast.Inspect(file, r.rewriteDeclarations)
	return astrewrite.Walk(file, r.rewriteIdents).(*ast.File)
}

func isRewrite(v *ast.Ident) bool {
	return strings.HasSuffix(v.Name, "_") && v.Obj != nil && v.Obj.Kind == ast.Con
}
