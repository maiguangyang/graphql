package model

import (
	"fmt"
	"unicode"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

var goTypeMap = map[string]string{
	"String"  :  "string",
	"Time"    :  "time.Time",
	"ID"      :  "string",
	"Float"   :  "float64",
	"Int"     :  "int64",
	"Boolean" :  "bool",
	"Any"     :  "interface{}",
}

type ObjectField struct {
	Def *ast.FieldDefinition
	Obj *Object
}

func (o *ObjectField) Name() string {
	return o.Def.Name.Value
}
func (o *ObjectField) MethodName() string {
	name := o.Name()
	return templates.ToGo(name)
}
func (o *ObjectField) LowerName() string {
  for i, v := range o.Name() {
    return string(unicode.ToUpper(v)) + o.Name()[i+1:]
  }
  return strcase.ToLowerCamel(o.Name())
}
func (o *ObjectField) Columns() []ObjectField {
	obj := o.TargetObject().Columns();
	return obj
}
func (o *ObjectField) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectField) IsColumn() bool {
	return o.HasDirective("column")
}
func (o *ObjectField) IsIdentifier() bool {
	return o.Name() == "id"
}
func (o *ObjectField) IsRelationship() bool {
	return o.HasDirective("relationship")
}
func (o *ObjectField) IsCreatable() bool {
	return !(o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}
func (o *ObjectField) IsUpdatable() bool {
	return !(o.Name() == "id" || o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}
func (o *ObjectField) IsReadonlyType() bool {
	return !(o.IsScalarType() || o.IsEnumType()) || o.Obj.Model.HasObject(o.TargetType())
}
func (o *ObjectField) IsWritableType() bool {
	return !o.IsReadonlyType()
}
func (o *ObjectField) IsScalarType() bool {
	return o.Obj.Model.HasScalar(o.TargetType())
}
func (o *ObjectField) IsEnumType() bool {
	return o.Obj.Model.HasEnum(o.TargetType())
}
func (o *ObjectField) IsOptional() bool {
	return !isNonNullType(o.Def.Type)
}
func (o *ObjectField) IsList() bool {
	return isListType(o.Def.Type)
}
func (o *ObjectField) IsEmbedded() bool {
	return !o.IsColumn() && !o.IsRelationship()
}
func (o *ObjectField) HasTargetObject() bool {
	return o.Obj.Model.HasObject(o.TargetType())
}
func (o *ObjectField) TargetObject() *Object {
	obj := o.Obj.Model.Object(o.TargetType())
	return &obj
}
func (o *ObjectField) HasTargetObjectExtension() bool {
	return o.Obj.Model.HasObjectExtension(o.TargetType())
}
func (o *ObjectField) TargetObjectExtension() *ObjectExtension {
	e := o.Obj.Model.ObjectExtension(o.TargetType())
	return &e
}
func (o *ObjectField) IsSortable() bool {
	return !o.IsReadonlyType() && o.IsScalarType()
}
func (o *ObjectField) IsSearchable() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "String" || t.Name.Value == "Int" || t.Name.Value == "Float"
}
func (o *ObjectField) IsString() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "String"
}
func (o *ObjectField) IsPassWord() bool {
	return o.Name() == "password"
}
func (o *ObjectField) IsState() bool {
	return o.Name() == "state"
}
func (o *ObjectField) IsDel() bool {
	return o.Name() == "del"
}
func (o *ObjectField) IsRequired() bool {
	return isNonNullType(o.Def.Type)
}
func (o *ObjectField) Directive(name string) *ast.Directive {
	for _, d := range o.Def.Directives {
		if d.Name.Value == name {
			return d
		}
	}
	return nil
}
func (o *ObjectField) NeedsQueryResolver() bool {
	return o.IsEmbedded()
}
func (o *ObjectField) HasDirective(name string) bool {
	return o.Directive(name) != nil
}
func (o *ObjectField) HasTargetTypeWithIDField() bool {
	if o.HasTargetObject() && o.TargetObject().HasField("id") {
		return true
	}
	if o.HasTargetObjectExtension() && o.TargetObjectExtension().Object.HasField("id") {
		return true
	}
	return false
}

func (o *ObjectField) GoType() string {
	return o.GoTypeWithPointer(true)
}
func (o *ObjectField) GoTypeWithPointer(showPointer bool) string {
	t := o.Def.Type
	st := ""

	if o.IsOptional() && showPointer {
		st += "*"
	} else {
		t = getNullableType(t)
	}

	if isListType(t) {
		st += "[]*"
	}

	v, ok := getNamedType(o.Def.Type).(*ast.Named)
	if ok {
		_t, known := goTypeMap[v.Name.Value]
		if known {
			st += _t
		} else {
			st += v.Name.Value
		}
	}

	return st
}

// maiguangyang new add
func (o *ObjectField) GetArgValue(name string) map[string]map[string]string {
	for _, d := range o.Def.Directives {
		if d.Name.Value == name && len(d.Arguments) > 0 {
			argArr := map[string]map[string]string{
				name: map[string]string{},
			}
			for _, child := range d.Arguments {
				argArr[name][child.Name.Value] = child.Value.GetValue().(string)
			}
			return argArr
		}
	}

	return map[string]map[string]string{}
}

// 获取字段说明
func (o *ObjectField) GetComment() string {
	column := o.GetArgValue("column")
	value  := column["column"]["gorm"]
	str    := ""
	if value != "" {
		str = RegexpReplace(value, `comment '`, `';`)
	} else {
    switch o.Name() {
      case "id":
        str = "uuid"
      case "createdAt":
        str = "创建时间"
      case "updatedAt":
        str = "更新时间"
      case "deletedBy":
        str = "删除人"
      case "updatedBy":
        str = "更新人"
      case "createdBy":
        str = "创建人"
      case "state":
        str = "状态：1/正常、2/禁用、3/下架"
      case "del":
        str = "状态：1/正常、2/删除"
    }
	}
	return str
}

// 备注说明字段
func (o *ObjectField) GetRemark() string {
	str := ""
  switch o.Name() {
    case "id":
      str = "create方法不是必填"
  }
	return str
}

// 获取字段说明
func (o *ObjectField) GetType() string {
	column := o.GetArgValue("column")
	value  := column["column"]["gorm"]
	str    := ""

	if value != "" {
		str = RegexpReplace(value, `type:`, ` `)
	} else {
    switch o.Name() {
      case "id":
        str = "varchar(36)"
      case "createdAt":
        str = "bigint(13)"
      case "updatedAt":
        str = "bigint(13)"
      case "deletedBy":
        str = "varchar(36)"
      case "updatedBy":
        str = "varchar(36)"
      case "createdBy":
        str = "varchar(36)"
      case "state":
        str = "int(2)"
      case "del":
        str = "int(2)"
    }
	}
	return str
}

// 获取正则验证
func (o *ObjectField) GetValidator() string {
	column := o.GetArgValue("validator")
	value  := column["validator"]["type"]
	str    := ""
	if value != "" {
		str = value
	} else {
    switch o.Name() {
      case "state":
        str = "justInt"
      case "del":
        str = "justInt"
    }
	}
	return str
}

// 获取Arguments
func (o *ObjectField) Arguments() string {
	argString := ""

	for key, child := range o.Obj.Def.Fields[0].Arguments {
		nullType  := ""
		bool := isNonNullType(child.Type)
		if bool {
			nullType = "!"
		}
		v, ok := getNamedType(child.Type).(*ast.Named)
		if ok {
			if key != len(o.Obj.Def.Fields[0].Arguments) - 1 {
				argString = argString + "$" + child.Name.Value + ": " + v.Name.Value + nullType + ", "
			} else {
				argString = argString + "$" + child.Name.Value + ": " + v.Name.Value + nullType
			}
		}
	}

	if argString != "" {
		argString = "(" + argString + ")"
	}
	return argString
}

// 获取Input
func (o *ObjectField) Inputs() string {
	argString := ""

	for key, child := range o.Obj.Def.Fields[0].Arguments {
		if key != len(o.Obj.Def.Fields[0].Arguments) - 1 {
			argString = argString + child.Name.Value + ": $" + child.Name.Value + ", "
		} else {
			argString = argString + child.Name.Value + ": $" + child.Name.Value
		}
	}

	if argString != "" {
		argString = "(" + argString + ")"
	}

	return argString
}

// 获取Field
func (o *ObjectField) Fields() string {
	argString := ""
	// { "name": "", "desc": "", "type": "", "required": "", "validator": "", "remark": "" },

	for key, child := range o.Obj.Def.Fields[0].Arguments {
		nullType  := "false"
		bool := isNonNullType(child.Type)
		if bool {
			nullType = "true"
		}
		v, ok := getNamedType(child.Type).(*ast.Named)
		if ok {
			if key != len(o.Obj.Def.Fields[0].Arguments) - 1 {
				argString = argString + `{ name: '` + child.Name.Value + `', desc: '` + child.Name.Value + `', type: '` + v.Name.Value + `', required: '` + nullType + `', validator: '', remark: '' },` + "\n      "
			} else {
				argString = argString + `{ name: '` + child.Name.Value + `', desc: '` + child.Name.Value + `', type: '` + v.Name.Value + `', required: '` + nullType + `', validator: '', remark: '' },` + "\n"
			}
		}
	}
	return argString
}


// 表名
func (o *ObjectField) EntityName() string {
	if len(o.Obj.Def.Directives) > 0 && len(o.Obj.Def.Directives[0].Arguments) > 0 {
		title := o.Obj.Def.Directives[0].Arguments[0].Value.GetValue()
		return title.(string)
	}
	return o.Name()
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


func (o *ObjectField) ModelTags() string {
	_gorm := fmt.Sprintf("default:null")
	_valid := ""

	dateArr := []interface{}{"createdAt", "updatedAt", "deletedAt", "state", "createdBy", "updatedBy", "deletedBy"}
	fields := []interface{}{"required", "type", "repeat", "edit"}

	if o.Name() == "id" {
		_gorm = "type:varchar(36) comment 'uuid';primary_key;NOT NULL;"
	}

	if IndexOf(dateArr, o.Name()) != -1 {
		tye := "type:varchar(255)"

    comment := "null;default:null"
    switch o.Name() {
      case "createdAt":
      	tye = "type:bigint(13)"
        comment = "'创建时间';default:null"
      case "updatedAt":
      	tye = "type:bigint(13)"
        comment = "'更新时间';default:null"
      case "deletedAt":
      	tye = "type:bigint(13)"
        comment = "'删除时间';default:null"
      case "createdBy":
      	tye = "type:varchar(36)"
        comment = "'创建人';default:null"
      case "updatedBy":
      	tye = "type:varchar(36)"
        comment = "'更新人';default:null"
      case "deletedBy":
      	tye = "type:varchar(36)"
        comment = "'删除人';default:null"
      case "state":
      	tye = "type:int(2)"
        comment = "'状态：1/正常、2/禁用、3/下架';NOT NULL;default:1;"
        _valid = "required:true;type:state;"
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
		} else if d.Name.Value == "validator" {
			for _, arg := range d.Arguments {
				if arg.Value.GetValue() != nil && IndexOf(fields, arg.Name.Value) != -1 {
					_valid += fmt.Sprintf("%v", arg.Name.Value + ":" + arg.Value.GetValue().(string) + ";")
				}
				// if arg.Name.Value == "required" && arg.Value.GetValue() != nil || arg.Name.Value == "type" && arg.Value.GetValue() != nil || arg.Name.Value == "repeat" && arg.Value.GetValue() != nil {
				// 	_valid += fmt.Sprintf("%v", arg.Name.Value + ":" + arg.Value.GetValue().(string) + ";")
				// }
			}
		}
	}

	str := fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), _gorm)

	if _valid != "" {
		str = fmt.Sprintf(`json:"%s" gorm:"%s" validator:"%s"`, o.Name(), _gorm, _valid)
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

func (o *ObjectField) FilterMapping() []FilterMappingItem {
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
