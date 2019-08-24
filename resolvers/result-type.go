package resolvers

import (
	// "fmt"
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"

	"github.com/jinzhu/gorm"
)

type EntityFilter interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, dialect gorm.Dialect, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntitySort interface {
	String() string
}

type EntityResultType struct {
	Offset       *int
	Limit        *int
	CurrentPage  *int
	PerPage      *int
	Query        EntityFilterQuery
	Sort         []EntitySort
	Filter       EntityFilter
	Fields       []*ast.Field
	SelectionSet *ast.SelectionSet
}

// maiguangyang new add
// 驼峰转蛇线
func snakeString(s string) string {
    data := make([]byte, 0, len(s)*2)
    j := false
    num := len(s)
    for i := 0; i < num; i++ {
        d := s[i]
        if i > 0 && d >= 'A' && d <= 'Z' && j {
            data = append(data, '_')
        }
        if d != '_' {
            j = true
        }
        data = append(data, d)
    }
    return strings.ToLower(string(data[:]))
}

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
				fields = append(fields, alias + "." + snakeString(sel.Name))
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

	dialect := q.Dialect()

	for _, s := range r.Sort {
		direction := "ASC"
		_s := s.String()
		if strings.HasSuffix(_s, "_DESC") {
			direction = "DESC"
		}
		col := strcase.ToLowerCamel(strings.ToLower(strings.TrimSuffix(_s, "_"+direction)))
		q = q.Order(dialect.Quote(col) + " " + direction)
	}

	wheres := []string{"state = ?"}
	values := []interface{}{1}
	joins := []string{}

	err := r.Query.Apply(ctx, dialect, r.SelectionSet, &wheres, &values, &joins)
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

	uniqueJoinsMap := map[string]bool{}
	uniqueJoins := []string{}
	for _, join := range joins {
		if !uniqueJoinsMap[join] {
			uniqueJoinsMap[join] = true
			uniqueJoins = append(uniqueJoins, join)
		}
	}

	for _, join := range uniqueJoins {
		q = q.Joins(join)
	}

	q = q.Group(alias + ".id")
	return q.Find(out).Error
}

// GetTotal ...
func (r *EntityResultType) GetTotal(ctx context.Context, db *gorm.DB, out interface{}) (count int, err error) {
	q := db

	dialect := q.Dialect()
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err = r.Query.Apply(ctx, dialect, r.SelectionSet, &wheres, &values, &joins)
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

	uniqueJoinsMap := map[string]bool{}
	uniqueJoins := []string{}
	for _, join := range joins {
		if !uniqueJoinsMap[join] {
			uniqueJoinsMap[join] = true
			uniqueJoins = append(uniqueJoins, join)
		}
	}

	for _, join := range uniqueJoins {
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
