package model

import (
	"fmt"
	// "strings"

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
	str := ""
	for _, d := range o.Def.Directives {
		if d.Name.Value == "validator" {
			for _, arg := range d.Arguments {
					v, ok := arg.Value.GetValue().(string)
					str += arg.Name.Value + ":" + v + ";"
					if !ok {
						panic(fmt.Sprintf("invalid value for %s->%s validator", o.Obj.Name(), o.Name()))
					}
			}
			return str
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

func (o *ObjectValid) ModelTags() string {
	valid := o.InverseValidatorName()
	return fmt.Sprintf(`json:"%s" validator:"%s"`, o.Name(), valid)
}

