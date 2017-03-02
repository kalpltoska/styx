package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/doctor-fate/styx/rewriter"
)

func main() {
	file, err := parser.ParseFile(token.NewFileSet(), os.Args[1], nil, 0)
	if err != nil {
		log.Print(err)
	}

	r := rewriter.NewIdentRewriter()
	rewritten := r.Rewrite(file)

	var buffer bytes.Buffer
	if err := format.Node(&buffer, token.NewFileSet(), rewritten); err != nil {
		log.Fatal(err)
	}
	fmt.Print(buffer.String())
}
