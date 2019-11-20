package model

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

func filterDirective(ds []*ast.Directive, name string) []*ast.Directive {
	res := []*ast.Directive{}
	for _, d := range ds {
		if d.Name.Value != name {
			res = append(res, d)
		}
	}
	return res
}

// PrintSchema
func PrintSchema(model Model) (string, error) {

	for _, o := range model.Objects() {
		fields := []*ast.FieldDefinition{}
		for _, f := range o.Def.Fields {
			f.Directives = filterDirective(f.Directives, "relationship")
			f.Directives = filterDirective(f.Directives, "column")
			f.Directives = filterDirective(f.Directives, "validator")
			fields = append(fields, f)
		}
		o.Def.Fields = fields
		o.Def.Directives = filterDirective(o.Def.Directives, "entity")
	}

	for _, o := range model.ObjectExtensions() {
		o.Object.Def.Directives = filterDirective(o.Object.Def.Directives, "entity")
	}

	printed := printer.Print(model.Doc)
	printedString, _ := printed.(string)

	return printedString, nil
}
