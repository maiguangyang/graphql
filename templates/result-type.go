package templates

var ResultType = `package gen

import (
	"context"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"

	"github.com/jinzhu/gorm"
)

func GetItem(ctx context.Context, db *gorm.DB, out interface{}, id *string) error {
	return db.Find(out, "id = ?", id).Error
}

func GetItemForRelation(ctx context.Context, db *gorm.DB, obj interface{}, relation string, out interface{}) error {
	return db.Model(obj).Related(out, relation).Error
}

type EntityFilter interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, dialect gorm.Dialect, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error
}

type EntitySort interface {
	Apply(ctx context.Context, dialect gorm.Dialect, sorts *[]string, joins *[]string) error
}

// type EntitySort interface {
// 	String() string
// }

type Created struct {
  name string
}

func (b Created) String() string {
  return "CREATED_AT_DESC"
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

      if len(sel.ObjectDefinition.Fields) > 0 {
  			for _, child := range sel.ObjectDefinition.Fields {
          // fmt.Println("child:", child.Name, child.Type.Elem == nil, child.Type.Elem, child.Type.Name(), child.Type.String())
          if child.Type.Elem == nil &&  child.Name != "assignee" && strings.Index(sel.Name, "Ids") == -1 {
            fields = append(fields, alias + "." + snakeString(child.Name))
          }
  			}
      }
			// ignore private field names !strings.HasPrefix(sel.Name, "__") &&
			if !strings.HasPrefix(sel.Name, "__") && len(sel.SelectionSet) == 0 && strings.Index(sel.Name, "Ids") == -1 {
				fields = append(fields, alias + "." + snakeString(sel.Name))
			}
      // else if (len(sel.SelectionSet) > 0) {
			// 	fields = append(fields, alias + "." + sel.Name)
			// }
		// case *ast.InlineFragment:
		// 	fields = recurseSelectionSets(reqCtx, fields, sel.SelectionSet)
		// case *ast.FragmentSpread:

		// 	fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
		// 	fields = recurseSelectionSets(reqCtx, fields, fragment.SelectionSet)
		}
	}
	return fields
}

// 查找数组并返回下标
// func indexOf(str []EntitySort, data string) int {
//   for k, v := range str{
//     if v.String() == data {
//       return k
//     }
//   }
//   return - 1
// }

// 查找数组并返回下标
func IndexOfTwo(str []string, data interface{}) int {
  for k, v := range str{
    if v == data {
      return k
    }
  }

  return - 1
}

type GetItemsOptions struct {
	Alias      string
	Preloaders []string
}

// GetResultTypeItems ...
func (r *EntityResultType) GetData(ctx context.Context, db *gorm.DB, opts GetItemsOptions, out interface{}) error {
	q := db

	// 麦广扬添加
	selects := GetFieldsRequested(ctx, opts.Alias)
	if len(selects) > 0 && IndexOfTwo(selects, opts.Alias + ".id") == -1 {
		selects = append(selects, opts.Alias + ".id")
	}

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


  // 以前是改写排序
 //  var _newSort []EntitySort
 //  // 判断是否有创建时间排序，没有的话，追加一个默认值
	// if indexOf(r.Sort, "CREATED_AT_DESC") == -1 && indexOf(r.Sort, "CREATED_AT_ASC") == -1 {
 //    var created Created
 //    var _sort EntitySort = created

 //    _newSort = append(_newSort, _sort)
	// }

 //  for _, v := range r.Sort {
 //    _newSort = append(_newSort, v)
 //  }
 //  r.Sort = _newSort

 //  // 原来的代码
	// for _, s := range r.Sort {
	// 	direction := "ASC"
	// 	_s := s.String()
	// 	if strings.HasSuffix(_s, "_DESC") {
	// 		direction = "DESC"
	// 	}

 //    // strcase.ToLowerCamel(
 //    col := strings.ToLower(strings.TrimSuffix(_s, "_"+direction))
 //    q = q.Order(opts.Alias + "." + dialect.Quote(col) + " " + direction)
	// }

	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	sorts := []string{}

	for _, sort := range r.Sort {
		sort.Apply(ctx, dialect, &sorts, &joins)
	}

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

	if len(sorts) > 0 {
		q = q.Order(strings.Join(sorts, ", "))
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

	if len(opts.Preloaders) > 0 {
		for _, p := range opts.Preloaders {
			q = q.Preload(p)
		}
	}
	q = q.Group(opts.Alias + ".id")
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

	if len(wheres) > 0 {
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
	err = q.Model(out).Count(&count).Error
	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
`
