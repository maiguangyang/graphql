package resolvers

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/ast"
	"github.com/iancoleman/strcase"

	"github.com/jinzhu/gorm"
)

type EntityFilter interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntitySort interface {
	String() string
}

type EntityResultType struct {
	Offset *int
	Limit  *int
	Query  EntityFilterQuery
	Sort   []EntitySort
	Filter EntityFilter
}

// maiguangyang new add
func GetFieldsRequested(ctx context.Context) []string {
	reqCtx := graphql.GetRequestContext(ctx)
	fieldSelections := graphql.GetResolverContext(ctx).Field.Selections
	return recurseSelectionSets(reqCtx, []string{}, fieldSelections)
}

// maiguangyang new add
func recurseSelectionSets(reqCtx *graphql.RequestContext, fields []string, selection ast.SelectionSet) []string {
	for _, sel := range selection {
		switch sel := sel.(type) {
		case *ast.Field:
			// ignore private field names
			if !strings.HasPrefix(sel.Name, "__") && len(sel.SelectionSet) == 0 {
				fields = append(fields, sel.Name)
			}
		// case *ast.InlineFragment:
		// 	fields = recurseSelectionSets(reqCtx, fields, sel.SelectionSet)
		// case *ast.FragmentSpread:

		// 	fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
		// 	fields = recurseSelectionSets(reqCtx, fields, fragment.SelectionSet)
		}
	}
	return fields
}


// GetResultTypeItems ...
func (r *EntityResultType) GetItems(ctx context.Context, db *gorm.DB, alias string, out interface{}) error {
	q := db

	// 麦广扬添加
	selects := GetFieldsRequested(ctx)
	if len(selects) > 0 {
		q = q.Select(selects)
	}

	// 原来的
	if r.Limit != nil {
		q = q.Limit(*r.Limit)
	}
	if r.Offset != nil {
		q = q.Offset(*r.Offset)
	}

	for _, s := range r.Sort {
		direction := "ASC"
		_s := s.String()
		if strings.HasSuffix(_s, "_DESC") {
			direction = "DESC"
		}
		col := strcase.ToLowerCamel(strings.ToLower(strings.TrimSuffix(_s, "_"+direction)))
		q = q.Order(col + " " + direction)
	}

	dialect := q.Dialect()
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err := r.Query.Apply(ctx, dialect, &wheres, &values, &joins)
	if err != nil {
		return err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &values, &joins)
		if err != nil {
			return err
		}
	}

	if len(wheres) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), values...)
	}

	uniqueJoins := map[string]bool{}
	for _, join := range joins {
		uniqueJoins[join] = true
	}

	for join := range uniqueJoins {
		q = q.Joins(join)
	}

	q = q.Group(alias + ".id")
	return q.Find(out).Error
}

// GetCount ...
func (r *EntityResultType) GetCount(ctx context.Context, db *gorm.DB, out interface{}) (count int, err error) {
	q := db

	dialect := q.Dialect()
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err = r.Query.Apply(ctx, dialect, &wheres, &values, &joins)
	if err != nil {
		return 0, err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &values, &joins)
		if err != nil {
			return 0, err
		}
	}

	if len(wheres) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), values...)
	}

	uniqueJoins := map[string]bool{}
	for _, join := range joins {
		uniqueJoins[join] = true
	}

	for join := range uniqueJoins {
		q = q.Joins(join)
	}
	err = q.Model(out).Count(&count).Error
	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
