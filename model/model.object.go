package model

import (
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

type Object struct {
	Def       *ast.ObjectDefinition
	Model     *Model
	Extension *ObjectExtension
}

func (o *Object) Name() string {
	return o.Def.Name.Value
}
func (o *Object) PluralName() string {
	return inflection.Plural(o.Name())
}
func (o *Object) LowerName() string {
	return strcase.ToLowerCamel(o.Def.Name.Value)
}
func (o *Object) TableName() string {
	return strcase.ToSnake(inflection.Plural(o.LowerName()))
}
func (o *Object) Column(name string) *ObjectColumn {
	for _, f := range o.Def.Fields {
		if o.isColumn(f) && f.Name.Value == name {
			return &ObjectColumn{f, o}
		}
	}
	return nil
}
func (o *Object) Columns() []ObjectColumn {
	columns := []ObjectColumn{}
	for _, f := range o.Def.Fields {
		if o.isColumn(f) {
			columns = append(columns, ObjectColumn{f, o})
		}
	}
	return columns
}
func (o *Object) HasReadonlyColumns() bool {
	for _, c := range o.Columns() {
		if c.IsReadonlyType() {
			return true
		}
	}
	return false
}
func (o *Object) IsToManyColumn(c ObjectColumn) bool {
	if c.Obj.Name() != o.Name() {
		return false
	}
	return o.HasRelationship(strings.TrimSuffix(c.Name(), "Ids"))
}
func (o *Object) Relationships() []*ObjectRelationship {
	relationships := []*ObjectRelationship{}
	for _, f := range o.Def.Fields {
		if o.isRelationship(f) {
			relationships = append(relationships, &ObjectRelationship{f, o})
		}
	}
	return relationships
}

func (o *Object) Relationship(name string) *ObjectRelationship {
	for _, rel := range o.Relationships() {
		if rel.Name() == name {
			return rel
		}
	}
	panic(fmt.Sprintf("relationship %s->%s not found", o.Name(), name))
}
func (o *Object) HasAnyRelationships() bool {
	return len(o.Relationships()) > 0
}
func (o *Object) HasRelationship(name string) bool {
	for _, rel := range o.Relationships() {
		if rel.Name() == name {
			return true
		}
	}
	return false
}
func (o *Object) Directive(name string) *ast.Directive {
	for _, d := range o.Def.Directives {
		if d.Name.Value == name {
			return d
		}
	}
	return nil
}
func (o *Object) HasDirective(name string) bool {
	return o.Directive(name) != nil
}

func (o *Object) isColumn(f *ast.FieldDefinition) bool {
	return !o.isRelationship(f)
}
func (o *Object) isRelationship(f *ast.FieldDefinition) bool {
	for _, d := range f.Directives {
		if d != nil && d.Name.Value == "relationship" {
			return true
		}
	}
	return false
}
func (o *Object) IsExtended() bool {
	return o.Extension != nil
}
