package templates

var GeneratedResolver = `package gen

import (
  "context"
  "time"
  "math"
  "strings"

  "github.com/99designs/gqlgen/graphql"
  "github.com/gofrs/uuid"
  "github.com/maiguangyang/graphql/events"
  "github.com/maiguangyang/graphql/resolvers"
  "github.com/vektah/gqlparser/ast"
)

func getPrincipalID(ctx context.Context) *string {
  v, _ := ctx.Value(KeyPrincipalID).(*string)
  return v
}

type GeneratedResolver struct {
  DB *DB
  EventController *events.EventController
}

func (r *GeneratedResolver) Mutation() MutationResolver {
  return &GeneratedMutationResolver{r}
}
func (r *GeneratedResolver) Query() QueryResolver {
  return &GeneratedQueryResolver{r}
}
{{range .Model.Objects}}
func (r *GeneratedResolver) {{.Name}}ResultType() {{.Name}}ResultTypeResolver {
  return &Generated{{.Name}}ResultTypeResolver{r}
}
{{if .HasAnyRelationships}}
func (r *GeneratedResolver) {{.Name}}() {{.Name}}Resolver {
  return &Generated{{.Name}}Resolver{r}
}
{{end}}
{{end}}

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range .Model.Objects}}
func (r *GeneratedMutationResolver) Create{{.Name}}(ctx context.Context, input map[string]interface{}) (item *{{.Name}}, err error) {
  principalID := getPrincipalID(ctx)
  now := time.Now()
  item = &{{.Name}}{ID: uuid.Must(uuid.NewV4()).String(), CreatedBy: principalID}
  tx := r.DB.db.Begin()

  event := events.NewEvent(events.EventMetadata{
    Type:        events.EventTypeCreated,
    Entity:      "{{.Name}}",
    EntityID:    item.ID,
    Date:        now.Unix(),
    PrincipalID: principalID,
  })

  var changes {{.Name}}Changes
  err = ApplyChanges(input, &changes)
  if err != nil {
    return
  }

{{range $col := .Columns}}{{if $col.IsCreatable}}
  if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
    item.{{$col.MethodName}} = changes.{{$col.MethodName}}
    event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
  }
{{end}}
{{end}}

  if err = tx.Create(item).Error; err != nil {
    return
  }

{{range $rel := .Relationships}}
{{if $rel.IsToMany}}
  items := []{{$rel.TargetType}}{}

  if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
    tx.Find(&items, "id IN (?)", ids)
    if err = tx.Model(&item).Association("{{$rel.MethodName}}").Replace(items).Error; err != nil {
      tx.Rollback()
      return
    }
  }

  if err = tx.Model(&items).Where("assigneeId = ?", item.ID).Update("state", item.State).Error; err != nil {
    tx.Rollback()
    return
  }

{{end}}
{{end}}

  // if err != nil {
  //  tx.Rollback()
  //  return
  // }
  err = tx.Commit().Error
  if err != nil {
    tx.Rollback()
    return
  }

  if len(event.Changes) > 0 {
    err = r.EventController.SendEvent(ctx, &event)
  }

  return
}
func (r *GeneratedMutationResolver) Update{{.Name}}(ctx context.Context, id string, input map[string]interface{}) (item *{{.Name}}, err error) {
  principalID := getPrincipalID(ctx)
  item = &{{.Name}}{}
  now := time.Now()
  tx := r.DB.db.Begin()

  event := events.NewEvent(events.EventMetadata{
    Type:        events.EventTypeUpdated,
    Entity:      "{{.Name}}",
    EntityID:    id,
    Date:        now.Unix(),
    PrincipalID: principalID,
  })

  var changes {{.Name}}Changes
  err = ApplyChanges(input, &changes)
  if err != nil {
    return
  }

  err = resolvers.GetItem(ctx, tx, item, &id)
  if err != nil {
    return
  }

  {{range $rel := .Relationships}}
  {{if $rel.IsToMany}}
    oldState       := item.State
  {{end}}
  {{end}}

  item.UpdatedBy = principalID

{{range $col := .Columns}}{{if $col.IsUpdatable}}
  if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
    event.AddOldValue("{{$col.Name}}", item.{{$col.MethodName}})
    event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
    item.{{$col.MethodName}} = changes.{{$col.MethodName}}
  }
{{end}}
{{end}}

  if err = tx.Save(item).Error; err != nil {
    return
  }

{{range $rel := .Relationships}}
{{if $rel.IsToMany}}
  items := []{{$rel.TargetType}}{}

  if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
    tx.Find(&items, "id IN (?)", ids)
    if err = tx.Model(&item).Association("{{$rel.MethodName}}").Replace(items).Error; err != nil {
      tx.Rollback()
      return
    }
  }

  // 判断是不是改变状态
  if oldState != item.State {
    if err = tx.Model(&items).Where("assigneeId = ?", item.ID).Update("state", item.State).Error; err != nil {
      tx.Rollback()
      return
    }
  }
{{end}}
{{end}}

  // if err != nil {
  //  tx.Rollback()
  //  return
  // }
  err = tx.Commit().Error
  if err != nil {
    tx.Rollback()
    return
  }

  if len(event.Changes) > 0 {
    err = r.EventController.SendEvent(ctx, &event)
    // data, _ := json.Marshal(event)
    // fmt.Println("?",string(data))
  }

  return
}
func (r *GeneratedMutationResolver) Delete{{.Name}}(ctx context.Context, id string) (item *{{.Name}}, err error) {
  principalID := getPrincipalID(ctx)
  item = &{{.Name}}{}
  now := time.Now()
  tx := r.DB.db.Begin()

  err = resolvers.GetItem(ctx, tx, item, &id)
  if err != nil {
    return
  }

  // 3为删除
  var state int64 = 3

  item.UpdatedBy  = principalID
  item.State      = &state

  // err = r.DB.Query().Delete(item, "{{.TableName}}.id = ?", id).Error

  event := events.NewEvent(events.EventMetadata{
    Type:        events.EventTypeDeleted,
    Entity:      "{{.Name}}",
    EntityID:    id,
    Date:        now.Unix(),
    PrincipalID: principalID,
  })

  if err = tx.Save(item).Error; err != nil {
    return
  }

{{range $rel := .Relationships}}
{{if $rel.IsToMany}}
  items := []{{$rel.TargetType}}{}
  // tx.Find(&items, "id = ?", id)
  // err = tx.Model(&items).Delete(items).Error
  if err = tx.Model(&items).Where("assigneeId = ?", id).Update("state", state).Error; err != nil {
    tx.Rollback()
    return
  }
{{end}}
{{end}}

  // if err != nil {
  //  tx.Rollback()
  //  return
  // }
  err = tx.Commit().Error
  if err != nil {
    tx.Rollback()
    return
  }
  err = r.EventController.SendEvent(ctx, &event)

  return
}

func (r *GeneratedMutationResolver) DeleteAll{{.PluralName}}(ctx context.Context) (bool, error) {
	err := r.DB.db.Delete(&{{.Name}}{}).Error
	return err == nil, err
}

{{end}}

type GeneratedQueryResolver struct{ *GeneratedResolver }

{{range $object := .Model.Objects}}
func (r *GeneratedQueryResolver) {{$object.Name}}(ctx context.Context, id *string, q *string, filter *{{$object.Name}}FilterType) (*{{$object.Name}}, error) {
  query := {{$object.Name}}QueryFilter{q}
  current_page := 0
  per_page := 0
  rt := &{{$object.Name}}ResultType{
    EntityResultType: resolvers.EntityResultType{
      CurrentPage: &current_page,
      PerPage:  &per_page,
      Query:  &query,
      Filter: filter,
    },
  }
  qb := r.DB.Query()
  if id != nil {
    qb = qb.Where("{{$object.TableName}}.id = ?", *id)
  }

  var items []*{{$object.Name}}
  err := rt.GetData(ctx, qb, "{{$object.TableName}}", &items)
  if err != nil {
    return nil, err
  }
  if len(items) == 0 {
    return nil, fmt.Errorf("{{$object.Name}} not found")
  }
  return items[0], err
}
func (r *GeneratedQueryResolver) {{$object.PluralName}}(ctx context.Context, current_page *int, per_page *int, q *string, sort []{{$object.Name}}SortType, filter *{{$object.Name}}FilterType) (*{{$object.Name}}ResultType, error) {
  _sort := []resolvers.EntitySort{}
  for _, s := range sort {
    _sort = append(_sort, s)
  }
  query := {{$object.Name}}QueryFilter{q}

  var selectionSet *ast.SelectionSet
  for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
    if f.Field.Name == "items" {
      selectionSet = &f.Field.SelectionSet
    }
  }

  return &{{$object.Name}}ResultType{
    EntityResultType: resolvers.EntityResultType{
      CurrentPage: current_page,
      PerPage:  per_page,
      Query:  &query,
      Sort: _sort,
      Filter: filter,
      SelectionSet: selectionSet,
    },
  }, nil
}

type Generated{{$object.Name}}ResultTypeResolver struct{ *GeneratedResolver }

func (r *Generated{{$object.Name}}ResultTypeResolver) Data(ctx context.Context, obj *{{$object.Name}}ResultType) (items []*{{$object.Name}}, err error) {
  err = obj.GetData(ctx, r.DB.db, "{{$object.TableName}}", &items)
  return
}

func (r *Generated{{$object.Name}}ResultTypeResolver) Total(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
  return obj.GetTotal(ctx, r.DB.db, &{{$object.Name}}{})
}

func (r *Generated{{$object.Name}}ResultTypeResolver) CurrentPage(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
  return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *Generated{{$object.Name}}ResultTypeResolver) PerPage(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
  return int(*obj.EntityResultType.PerPage), nil
}

func (r *Generated{{$object.Name}}ResultTypeResolver) TotalPage(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
  total, _   := r.Total(ctx, obj)
  perPage, _ := r.PerPage(ctx, obj)
  totalPage  := int(math.Ceil(float64(total) / float64(perPage)))

  return totalPage, nil
}

{{if .HasAnyRelationships}}
type Generated{{$object.Name}}Resolver struct { *GeneratedResolver }

{{range $index, $relationship := .Relationships}}
func (r *Generated{{$object.Name}}Resolver) {{$relationship.MethodName}}(ctx context.Context, obj *{{$object.Name}}) (res {{.ReturnType}}, err error) {
{{if $relationship.IsToMany}}
  selects := resolvers.GetFieldsRequested(ctx, strings.ToLower("{{$relationship.MethodName}}"))

  items := []*{{.TargetType}}{}
  err = r.DB.Query().Where("state = ?", 1).Select(selects).Model(obj).Related(&items, "{{$relationship.MethodName}}").Error
  res = items
{{else}}
  loaders := ctx.Value("loaders").(map[string]*dataloader.Loader)
  if obj.{{$relationship.MethodName}}ID != nil {
    item, _err := loaders["{{$relationship.Target.Name}}"].Load(ctx, dataloader.StringKey(*obj.{{$relationship.MethodName}}ID))()
    res = item.({{.ReturnType}})
    err = _err
  }
{{end}}
  return
}

{{if $relationship.IsToMany}}
func (r *Generated{{$object.Name}}Resolver) {{$relationship.MethodName}}Ids(ctx context.Context, obj *{{$object.Name}}) (ids []string, err error) {
  ids = []string{}
  items := []*{{$relationship.TargetType}}{}
  err = r.DB.Query().Model(obj).Select("{{$relationship.Target.TableName}}.id").Related(&items, "{{$relationship.MethodName}}").Error
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
