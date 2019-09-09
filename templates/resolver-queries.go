package templates

var ResolverQueries = `package gen

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/maiguangyang/graphql/events"
	"github.com/maiguangyang/graphql/resolvers"
	"github.com/vektah/gqlparser/ast"
)

type GeneratedQueryResolver struct{ *GeneratedResolver }

{{range $obj := .Model.Objects}}
	type Query{{$obj.Name}}HandlerOptions struct {
		ID *string
		Q      *string
		Filter *{{$obj.Name}}FilterType
	}
	func (r *GeneratedQueryResolver) {{$obj.Name}}(ctx context.Context, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error) {
		opts := Query{{$obj.Name}}HandlerOptions{
			ID: id,
			Q: q,
			Filter: filter,
		}
		return r.Handlers.Query{{$obj.Name}}(ctx, r.GeneratedResolver, opts)
	}
	func Query{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, opts Query{{$obj.Name}}HandlerOptions) (*{{$obj.Name}}, error) {
		query := {{$obj.Name}}QueryFilter{opts.Q}
	  current_page := 0
	  per_page := 0
		rt := &{{$obj.Name}}ResultType{
			EntityResultType: resolvers.EntityResultType{
      	CurrentPage: &current_page,
	      PerPage:  &per_page,
	      Query:  &query,
	      Filter: opts.Filter,
			},
		}
		qb := r.DB.Query()
		if opts.ID != nil {
			qb = qb.Where("{{$obj.TableName}}.id = ?", *opts.ID)
		}

		var items []*{{$obj.Name}}
		err := rt.GetData(ctx, qb, "{{$obj.TableName}}", &items)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			return nil, fmt.Errorf("{{$obj.Name}} not found")
		}
		return items[0], err
	}

	type Query{{$obj.PluralName}}HandlerOptions struct {
		CurrentPage *int
		PerPage  *int
		Q      *string
		Sort   []{{$obj.Name}}SortType
		Filter *{{$obj.Name}}FilterType
	}
	func (r *GeneratedQueryResolver) {{$obj.PluralName}}(ctx context.Context, current_page *int, per_page *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error) {
		opts := Query{{$obj.PluralName}}HandlerOptions{
      CurrentPage: current_page,
      PerPage:  per_page,
			Q: q,
			Sort: sort,
			Filter: filter,
		}
		return r.Handlers.Query{{$obj.PluralName}}(ctx, r.GeneratedResolver, opts)
	}
	func Query{{$obj.PluralName}}Handler(ctx context.Context, r *GeneratedResolver, opts Query{{$obj.PluralName}}HandlerOptions) (*{{$obj.Name}}ResultType, error) {
		_sort := []resolvers.EntitySort{}
		for _, s := range opts.Sort {
			_sort = append(_sort, s)
		}
		query := {{$obj.Name}}QueryFilter{opts.Q}

		var selectionSet *ast.SelectionSet
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}

		return &{{$obj.Name}}ResultType{
			EntityResultType: resolvers.EntityResultType{
				CurrentPage: opts.CurrentPage,
				PerPage:  opts.PerPage,
				Query:  &query,
				Sort: _sort,
				Filter: opts.Filter,
				SelectionSet: selectionSet,
			},
		}, nil
	}

	type Generated{{$obj.Name}}ResultTypeResolver struct{ *GeneratedResolver }

	func (r *Generated{{$obj.Name}}ResultTypeResolver) Data(ctx context.Context, obj *{{$obj.Name}}ResultType) (items []*{{$obj.Name}}, err error) {
	  err = obj.GetData(ctx, r.DB.db, "{{$obj.TableName}}", &items)
	  return
	}

	func (r *Generated{{$obj.Name}}ResultTypeResolver) Total(ctx context.Context, obj *{{$obj.Name}}ResultType) (count int, err error) {
	  return obj.GetTotal(ctx, r.DB.db, &{{$obj.Name}}{})
	}

	func (r *Generated{{$obj.Name}}ResultTypeResolver) CurrentPage(ctx context.Context, obj *{{$obj.Name}}ResultType) (count int, err error) {
	  return int(*obj.EntityResultType.CurrentPage), nil
	}

	func (r *Generated{{$obj.Name}}ResultTypeResolver) PerPage(ctx context.Context, obj *{{$obj.Name}}ResultType) (count int, err error) {
	  return int(*obj.EntityResultType.PerPage), nil
	}

	func (r *Generated{{$obj.Name}}ResultTypeResolver) TotalPage(ctx context.Context, obj *{{$obj.Name}}ResultType) (count int, err error) {
	  total, _   := r.Total(ctx, obj)
	  perPage, _ := r.PerPage(ctx, obj)
	  totalPage  := int(math.Ceil(float64(total) / float64(perPage)))

	  return totalPage, nil
	}

	{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}
		type Generated{{$obj.Name}}Resolver struct { *GeneratedResolver }

		{{range $col := $obj.Columns}}
			{{if $col.IsReadonlyType}}
			func (r *Generated{{$obj.Name}}Resolver) {{$col.MethodName}}(ctx context.Context, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error) {
				return r.Handlers.{{$obj.Name}}{{$col.MethodName}}(ctx, r, obj)
			}
			func {{$obj.Name}}{{$col.MethodName}}Handler(ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error) {
				{{if and (not $col.IsList) $col.HasTargetTypeWithIDField ($obj.HasColumn (print $col.Name "Id"))}}
				if obj.{{$col.MethodName}}ID != nil {
					res = &{{$col.TargetType}}{ID: *obj.{{$col.MethodName}}ID}
				}
				{{else}}
				err = fmt.Errorf("Resolver handler for {{$obj.Name}}{{$col.MethodName}} not implemented")
				{{end}}
				return
			}
			{{end}}
		{{end}}

		{{range $index, $rel := $obj.Relationships}}
			func (r *Generated{{$obj.Name}}Resolver) {{$rel.MethodName}}(ctx context.Context, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error) {
				return r.Handlers.{{$obj.Name}}{{$rel.MethodName}}(ctx, r, obj)
			}
			func {{$obj.Name}}{{$rel.MethodName}}Handler(ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error) {
				{{if $rel.IsToMany}}
					selects := resolvers.GetFieldsRequested(ctx, strings.ToLower("{{$rel.Target.TableName}}"))

					items := []*{{$rel.TargetType}}{}
					err = r.DB.Query().Select(selects).Model(obj).Related(&items, "{{$rel.MethodName}}").Error
					res = items
				{{else}}
					loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
					if obj.{{$rel.MethodName}}ID != nil {
						item, _err := loaders["{{$rel.Target.Name}}"].Load(ctx, dataloader.StringKey(*obj.{{$rel.MethodName}}ID))()
						res, _ = item.({{$rel.ReturnType}})
						{{if $rel.IsNonNull}}
						if res == nil {
							_err = fmt.Errorf("{{$rel.Target.Name}} with id '%s' not found",*obj.{{$rel.MethodName}}ID)
						}{{end}}
						err = _err
					}
				{{end}}
				return
			}
			{{if $rel.IsToMany}}
				func (r *Generated{{$obj.Name}}Resolver) {{$rel.MethodName}}Ids(ctx context.Context, obj *{{$obj.Name}}) (ids []string, err error) {
					ids = []string{}

					items := []*{{$rel.TargetType}}{}
					err = r.DB.Query().Model(obj).Select("{{$rel.Target.TableName}}.id").Related(&items, "{{$rel.MethodName}}").Error

					for _, item := range items {
						ids = append(ids, item.ID)
					}

					return
				}
			{{end}}

		{{end}}

	{{end}}

{{end}}
`
