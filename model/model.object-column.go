package model

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

var goTypeMap = map[string]string{
	"String":  "string",
	"Time":    "time.Time",
	"ID":      "string",
	"Float":   "float64",
	"Int":     "int64",
	"Boolean": "bool",
}

type ObjectColumn struct {
	Def *ast.FieldDefinition
	Obj *Object
}

func (o *ObjectColumn) Name() string {
	return o.Def.Name.Value
}
func (o *ObjectColumn) MethodName() string {
	name := o.Name()
	if name == "id" {
		return "ID"
	}
	if name == "uid" {
		return "UID"
	}
	if name == "url" {
		return "URL"
	}

	if strings.HasSuffix(name, "Id") {
		name = strings.TrimSuffix(name, "Id") + "ID"
	}

	return strcase.ToCamel(name)
}

func (o *ObjectColumn) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectColumn) IsCreatable() bool {
	return !(o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy")
}
func (o *ObjectColumn) IsUpdatable() bool {
	return !(o.Name() == "id" || o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy")
}
func (o *ObjectColumn) IsOptional() bool {
	return o.Def.Type.GetKind() != "NonNull"
}
func (o *ObjectColumn) IsSearchable() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "String"
}
func (o *ObjectColumn) GoType() string {
	return o.GoTypeWithPointer(true)
}
func (o *ObjectColumn) GoTypeWithPointer(showPointer bool) string {
	t := ""

	if o.IsOptional() && showPointer {
		t += "*"
	}

	v, ok := getNamedType(o.Def.Type).(*ast.Named)
	if ok {
		_t, known := goTypeMap[v.Name.Value]
		if known {
			t += _t
		} else {
			t += v.Name.Value
		}
	}
	return t
}

func (o *ObjectColumn) InverseValidatorName() map[string]string {
	str := map[string]string{}
	for _, d := range o.Def.Directives {
		if d.Name.Value == "column" || d.Name.Value == "validator" {
			str[d.Name.Value] = ""
			for _, arg := range d.Arguments {
					v, ok := arg.Value.GetValue().(string)
					if v != "" {
            if arg.Name.Value == "length" {
              str["gorm"] += `type:varchar(` + v + `) comment `
            } else if arg.Name.Value == "comment" {
              str["gorm"] += `'` + v + `';`
            } else if arg.Name.Value == "isNull" && v == "false" {
              str["gorm"] += `NOT NULL;`
            } else if arg.Name.Value == "value" {
              str["gorm"] += `default:`+ v +`;`
            } else if arg.Name.Value == "required" && v == "true" {
              str["validator"] += `required:`+ v +`;`
            } else if arg.Name.Value == "valid" {
              str["validator"] += `type:`+ v +`;`
            }

        //     else {
						  // str[d.Name.Value] += arg.Name.Value + ":" + v + ";"
        //     }
						// str += arg.Name.Value + ":" + v + ";"
					}
					if !ok {
						panic(fmt.Sprintf("invalid value for %s->%s validator", o.Obj.Name(), o.Name()))
					}
			}
		}
	}

	return str

	// panic(fmt.Sprintf("missing validator directive argument for %s->%s validator", o.Obj.Name(), o.Name()))
}

func (o *ObjectColumn) ModelTags() string {
	_gorm := fmt.Sprintf("column:%s", o.Name())
	if o.Name() == "id" {
		_gorm += ";primary_key"
	}

  valid := o.InverseValidatorName()

  str := fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), _gorm)

  if len(valid["gorm"]) > 0 {
    str = fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), valid["gorm"])
  }

  if len(valid["validator"]) > 0 {
    str = fmt.Sprintf(`json:"%s" gorm:"%s" validator:"%s"`, o.Name(), valid["gorm"], valid["validator"])
  }

  return str
}

type FilterMappingItem struct {
	Suffix      string
	Operator    string
	InputType   ast.Type
	ValueFormat string
}

func (f *FilterMappingItem) SuffixCamel() string {
	return strcase.ToCamel(f.Suffix)
}
func (f *FilterMappingItem) WrapValueVariable(v string) string {
	return fmt.Sprintf(f.ValueFormat, v)
}

func (o *ObjectColumn) FilterMapping() []FilterMappingItem {
	t := getNamedType(o.Def.Type)
	mapping := []FilterMappingItem{
		FilterMappingItem{"", "= ?", t, "%s"},
		FilterMappingItem{"_ne", "!= ?", t, "%s"},
		FilterMappingItem{"_gt", "> ?", t, "%s"},
		FilterMappingItem{"_lt", "< ?", t, "%s"},
		FilterMappingItem{"_gte", ">= ?", t, "%s"},
		FilterMappingItem{"_lte", "<= ?", t, "%s"},
		FilterMappingItem{"_in", "IN (?)", listType(nonNull(t)), "%s"},
	}
	_t := t.(*ast.Named)
	if _t.Name.Value == "String" {
		mapping = append(mapping,
			FilterMappingItem{"_like", "LIKE ?", t, "strings.Replace(strings.Replace(*%s,\"?\",\"_\",-1),\"*\",\"%%\",-1)"},
			FilterMappingItem{"_prefix", "LIKE ?", t, "fmt.Sprintf(\"%%s%%%%\",*%s)"},
			FilterMappingItem{"_suffix", "LIKE ?", t, "fmt.Sprintf(\"%%%%%%s\",*%s)"},
		)
	}
	return mapping
}
