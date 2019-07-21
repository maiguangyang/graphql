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
	CurrentPage *int
	PerPage  *int
	Query  EntityFilterQuery
	Sort   []EntitySort
	Filter EntityFilter
}

// maiguangyang new add
func GetFieldsRequested(ctx context.Context, alias string) []string {
	reqCtx := graphql.GetRequestContext(ctx)
	fieldSelections := graphql.GetResolverContext(ctx).Field.Selections
	return recurseSelectionSets(reqCtx, []string{}, fieldSelections, alias)
}

// maiguangyang new add
func recurseSelectionSets(reqCtx *graphql.RequestContext, fields []string, selection ast.SelectionSet, alias string) []string {
	for _, sel := range selection {
		switch sel := sel.(type) {
		case *ast.Field:
			// ignore private field names
			if !strings.HasPrefix(sel.Name, "__") && len(sel.SelectionSet) == 0 {
				fields = append(fields, alias + "." + sel.Name)
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
func (r *EntityResultType) GetData(ctx context.Context, db *gorm.DB, alias string, out interface{}) error {
	q := db

	// 麦广扬添加
	selects := GetFieldsRequested(ctx, alias)

	if len(selects) > 0 {
		q = q.Select(selects)
	}

	// maiguangyang update
	if r.PerPage != nil {
		if int(*r.PerPage) == 0 {
			q = q.Limit(1)
		} else {
			q = q.Limit(*r.PerPage)
		}
	}
	if r.CurrentPage != nil {
		// q = q.Offset(*r.CurrentPage)
		q = q.Offset((int(*r.CurrentPage) - 1) * int(*r.PerPage))
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

	if len(wheres) > 0 && len(values) > 0 {
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

// GetTotal ...
func (r *EntityResultType) GetTotal(ctx context.Context, db *gorm.DB, out interface{}) (count int, err error) {
	q := db

	// if r.Limit != nil {
	// 	q = q.Limit(*r.Limit)
	// }
	// if r.Offset != nil {
	// 	q = q.Offset(*r.Offset)
	// }

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

	if len(wheres) > 0 && len(values) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), values...)
	}

	uniqueJoins := map[string]bool{}
	for _, join := range joins {
		uniqueJoins[join] = true
	}

	for join := range uniqueJoins {
		q = q.Joins(join)
	}

	if err := q.Model(out).Count(&count).Error; err != nil {
		return 0, nil
	}

	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
