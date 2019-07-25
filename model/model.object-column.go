package model

import (
	"fmt"

	"github.com/99designs/gqlgen/codegen/templates"
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
	return templates.ToGo(name)
}

func (o *ObjectColumn) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectColumn) IsCreatable() bool {
	return !(o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "deletedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy" || o.Name() == "deletedBy")
}
func (o *ObjectColumn) IsUpdatable() bool {
	return !(o.Name() == "id" || o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "deletedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy" || o.Name() == "deletedBy")
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


// 查找数组并返回下标
func IndexOf(str []interface{}, data interface{}) int {
  for k, v := range str{
    if v == data {
      return k
    }
  }

  return - 1
}

func (o *ObjectColumn) ModelTags() string {
	_gorm := fmt.Sprintf("column:%s;null;default:null", o.Name())
	dateArr := []interface{}{"createdAt", "updatedAt", "state"}

	if o.Name() == "id" {
		_gorm = "type:varchar(36) comment 'uuid';primary_key;NOT NULL;"
	}

	if IndexOf(dateArr, o.Name()) != -1 {
		tye := "type:int(11)"

    comment := "null;default:null"
    switch o.Name() {
      case "createdAt":
        comment = "'创建时间';null;default:null"
      case "updatedAt":
        comment = "'更新时间';null;default:null"
      case "state":
      	tye = "type:int(2)"
        comment = "'状态：1/正常、2/禁用、3/删除';NOT NULL;default:1;"
    }

    _gorm = fmt.Sprintf("%s comment %s", tye, comment)

	}

	for _, d := range o.Def.Directives {
		if d.Name.Value == "column" {
			for _, arg := range d.Arguments {
				if arg.Name.Value == "gorm" {
					_gorm = fmt.Sprintf("%v", arg.Value.GetValue())
				}
			}
		}
	}

	return fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), _gorm)
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
