package model

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

type ObjectValid struct {
	Def *ast.FieldDefinition
	Obj *Object
}

func (o *ObjectValid) Name() string {
	return o.Def.Name.Value
}

func (o *ObjectValid) MethodName() string {
	name := o.Name()
	return strcase.ToCamel(name)
}

func (o *ObjectValid) InverseValidatorName() string {
	for _, d := range o.Def.Directives {
		if d.Name.Value == "validator" {
			for _, arg := range d.Arguments {
				if arg.Name.Value == "valid" || arg.Name.Value == "required" {
					v, ok := arg.Value.GetValue().(string)
					fmt.Println(v)
					if !ok {
						panic(fmt.Sprintf("invalid value for %s->%s validator", o.Obj.Name(), o.Name()))
					}
					// return v
				}
			}
		}
	}
	panic(fmt.Sprintf("missing validator directive argument for %s->%s validator", o.Obj.Name(), o.Name()))
}

func (o *ObjectValid) ReturnType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)

	return fmt.Sprintf("*%s", nt.Name.Value)
}

func (o *ObjectValid) GoType() string {
	return o.ReturnType()
}

func (o *ObjectValid) InverseValidator() *ObjectValid {
	return o.Obj.Validator(o.InverseValidatorName())
}

func (o *ObjectValid) ModelTags() string {
	invrel := o.InverseValidator()
	return fmt.Sprintf(`json:"%s" validator:"valid:%s"`, o.Name(), strings.ToLower(invrel.MethodName()))
}

