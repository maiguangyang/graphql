package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func objectResultTypeDefinition(o *Object) *ast.ObjectDefinition {
	return &ast.ObjectDefinition{
		Kind: kinds.ObjectDefinition,
		Name: nameNode(o.Name() + "ResultType"),
		Fields: []*ast.FieldDefinition{
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("data"),
				Type: nonNull(&ast.List{
					Kind: kinds.List,
					Type: nonNull(namedType(o.Name())),
				}),
			},

      // &ast.FieldDefinition{
      //  Kind: kinds.FieldDefinition,
      //  Name: nameNode("pages"),
      //  Type: nonNull(namedType("Any")),
      // },
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("total"),
				Type: nonNull(namedType("Int")),
			},
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("current_page"),
				Type: nonNull(namedType("Int")),
			},
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("per_page"),
				Type: nonNull(namedType("Int")),
			},
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("total_page"),
				Type: nonNull(namedType("Int")),
			},
		},
	}
}
