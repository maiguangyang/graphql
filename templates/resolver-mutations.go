package templates

var ResolverMutations = `package gen

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/maiguangyang/graphql/events"
	"github.com/vektah/gqlparser/ast"
)

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range $obj := .Model.ObjectEntities}}
	func (r *GeneratedMutationResolver) Create{{$obj.Name}}(ctx context.Context, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
		return r.Handlers.Create{{$obj.Name}}(ctx, r.GeneratedResolver, input)
	}
	func Create{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
		principalID := GetPrincipalIDFromContext(ctx)
		now := time.Now()
		item = &{{$obj.Name}}{ID: uuid.Must(uuid.NewV4()).String(), CreatedBy: principalID}
		tx := r.DB.db.Begin()
    defer func() {
      if r := recover(); r != nil {
        tx.Rollback()
      }
    }()

		event := events.NewEvent(events.EventMetadata{
			Type:        events.EventTypeCreated,
			Entity:      "{{$obj.Name}}",
			EntityID:    item.ID,
			Date:        now.UnixNano() / 1e6,
			PrincipalID: principalID,
		})

	  {{range $col := .Columns}}
	  {{if $col.IsState}}
		  if input["state"] == nil {
		    input["state"] = 1
		  }
	  {{end}}
	  {{end}}

		var changes {{$obj.Name}}Changes
		err = ApplyChanges(input, &changes)
		if err != nil {
			return
		}

		{{range $col := .Columns}}{{if $col.IsCreatable}}
			if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
				item.{{$col.MethodName}} = changes.{{$col.MethodName}}
				{{if $col.IsIdentifier}}event.EntityID = item.{{$col.MethodName}}{{end}}
				event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
			}
		{{end}}{{end}}

	  errText, repeat, resErr := utils.Validator(item, "create")
	  if resErr != nil {
	    return item, &errText
	  }

    // 这里处理是否允许字段重复存在
    if len(repeat) > 0 {
      whereData := ""
      whereValue :=  []interface{}{}
      for _, v := range repeat {
        whereData = utils.IsWhereEmpty(whereData, v + " = ?")
        whereValue = append(whereValue, input[v])
      }

      first := &{{$obj.Name}}{}
      if err := tx.Where(whereData, whereValue...).First(&first).Error; err == nil {
        return item, fmt.Errorf("exist")
      }
    }

	  {{range $col := .Columns}}
	  {{if $col.IsPassWord}}
	    if input["password"] != nil {
	      item.Password = utils.EncryptPassword(item.Password)
	    }
	  {{end}}
	  {{end}}

	  if err := tx.Create(item).Error; err != nil {
	  	tx.Rollback()
	    return item, err
	  }

		{{range $rel := $obj.Relationships}}
			{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
				{{$rel.Name}} := []{{$rel.TargetType}}{}
				if ids,exists:=input["{{$rel.Name}}Ids"]; exists {
					tx.Find(&{{$rel.Name}}, "id IN (?)", ids)
					association := tx.Model(&item).Association("{{$rel.MethodName}}")
					association.Replace({{$rel.Name}})
				}
				{{if $rel.IsOneToMany}}
        if err := tx.Model(&{{$rel.Name}}).Where("{{$rel.InverseRelationshipName}}_id = ?", item.ID).Update("state", item.State).Error; err != nil {
          tx.Rollback()
          return item, err
        }
        {{end}}
			{{end}}{{end}}
		{{end}}

		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return
		}

		if len(event.Changes) > 0 {
			err = r.Handlers.OnEvent(ctx, r, &event)
			if err != nil {
				return
			}
			err = r.EventController.SendEvent(ctx, &event)
		}

		return
	}
	func (r *GeneratedMutationResolver) Update{{$obj.Name}}(ctx context.Context, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
		return r.Handlers.Update{{$obj.Name}}(ctx, r.GeneratedResolver, id, input)
	}
	func Update{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
		principalID := GetPrincipalIDFromContext(ctx)
		item = &{{$obj.Name}}{}
		now := time.Now()
		tx := r.DB.db.Begin()
    defer func() {
      if r := recover(); r != nil {
        tx.Rollback()
      }
    }()

		event := events.NewEvent(events.EventMetadata{
			Type:        events.EventTypeUpdated,
			Entity:      "{{$obj.Name}}",
			EntityID:    id,
			Date:        now.UnixNano() / 1e6,
			PrincipalID: principalID,
		})

	  {{range $col := .Columns}}
	  {{if $col.IsState}}
		  if input["state"] == nil {
		    input["state"] = 1
		  }
	  {{end}}
	  {{end}}

		var changes {{$obj.Name}}Changes
		err = ApplyChanges(input, &changes)
		if err != nil {
			return
		}

		{{range $col := .Columns}}{{if $col.IsUpdatable}}
			if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
				event.AddOldValue("{{$col.Name}}", item.{{$col.MethodName}})
				event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
				item.{{$col.MethodName}} = changes.{{$col.MethodName}}
			}
		{{end}}
		{{end}}

	  errText, repeat, resErr := utils.Validator(item, "update")
	  if resErr != nil {
	    return item, &errText
	  }

    // 这里处理是否允许字段重复存在
    if len(repeat) > 0 {
      whereData := "id <> ?"
      whereValue :=  []interface{}{id}
      for _, v := range repeat {
        whereData = utils.IsWhereEmpty(whereData, v + " = ?")
        whereValue = append(whereValue, input[v])
      }

      first := &{{$obj.Name}}{}
      if err := tx.Where(whereData, whereValue...).First(&first).Error; err == nil {
        return item, fmt.Errorf("exist")
      }
    }

	  item.UpdatedBy = principalID
	  item.ID        = id

	  {{range $col := .Columns}}
	  {{if $col.IsPassWord}}
	    if input["password"] != nil {
	      item.Password = utils.EncryptPassword(item.Password)
	    }
	  {{end}}
	  {{end}}

	  if err := tx.Model(&item).Updates(item).Error; err != nil {
	  	tx.Rollback()
	    return item, err
	  }

		{{range $rel := $obj.Relationships}}
		{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
			{{$rel.Name}} := []{{$rel.TargetType}}{}

			if ids,exists:=input["{{$rel.Name}}Ids"]; exists {
				tx.Find(&{{$rel.Name}}, "id IN (?)", ids)

				association := tx.Model(&item).Association("{{$rel.MethodName}}")
				association.Replace({{$rel.Name}})
			}

      {{if $rel.IsOneToMany}}
      if err := tx.Model(&{{$rel.Name}}).Where("{{$rel.InverseRelationshipName}}_id = ?", item.ID).Update("state", item.State).Error; err != nil {
        tx.Rollback()
        return item, err
      }
      {{end}}
		{{end}}{{end}}
		{{end}}

    if err := GetItem(ctx, tx, item, &id); err != nil {
      tx.Rollback()
      return item, err
    }

		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return
		}

		if len(event.Changes) > 0 {
			err = r.Handlers.OnEvent(ctx, r, &event)
			if err != nil {
				return
			}
			err = r.EventController.SendEvent(ctx, &event)
			// data, _ := json.Marshal(event)
			// fmt.Println("?",string(data))
		}

		return
	}
	func (r *GeneratedMutationResolver) Delete{{$obj.Name}}(ctx context.Context, id string) (item *{{$obj.Name}}, err error) {
		return r.Handlers.Delete{{$obj.Name}}(ctx, r.GeneratedResolver, id)
	}
	func Delete{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, id string) (item *{{$obj.Name}}, err error) {
		principalID := GetPrincipalIDFromContext(ctx)
		item = &{{$obj.Name}}{}
		now := time.Now()
		tx := r.DB.db.Begin()
    defer func() {
      if r := recover(); r != nil {
        tx.Rollback()
      }
    }()

		err = GetItem(ctx, tx, item, &id)
		if err != nil {
			return
		}

		deletedAt      := now.UnixNano() / 1e6

		item.DeletedAt = &deletedAt
		item.DeletedBy = principalID
		item.UpdatedBy = principalID

		event := events.NewEvent(events.EventMetadata{
			Type:        events.EventTypeDeleted,
			Entity:      "{{$obj.Name}}",
			EntityID:    id,
			Date:        deletedAt,
			PrincipalID: principalID,
		})

		// err = tx.Where(TableName("{{$obj.TableName}}") + ".id = ?", id).Delete(item).Error
		// if err != nil {
		// 	tx.Rollback()
		// 	return
		// }

	  if err := tx.Save(item).Error; err != nil {
	  	tx.Rollback()
	    return item, err
	  }

		{{range $rel := .Relationships}}
		{{if $rel.IsToMany}}
		  {{$rel.Name}} := []{{$rel.TargetType}}{}
		  // if err := tx.Where("{{$rel.InverseRelationshipName}}_id = ?", id).Delete(&{{$rel.Name}}).Error; err != nil {
		  if err := tx.Model(&{{$rel.Name}}).Where("{{$rel.InverseRelationshipName}}_id = ?", id).Updates(&{{$rel.TargetType}}{DeletedAt: &deletedAt, UpdatedBy: principalID, DeletedBy: principalID}).Error; err != nil {
		    tx.Rollback()
		    return item, err
		  }
		{{end}}
		{{end}}

		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return
		}

		err = r.Handlers.OnEvent(ctx, r, &event)
		if err != nil {
			return
		}
		err = r.EventController.SendEvent(ctx, &event)

		return
	}
	func (r *GeneratedMutationResolver) DeleteAll{{$obj.PluralName}}(ctx context.Context) (bool, error) {
		return r.Handlers.DeleteAll{{$obj.PluralName}}(ctx, r.GeneratedResolver)
	}
	func DeleteAll{{$obj.PluralName}}Handler(ctx context.Context, r *GeneratedResolver) (bool, error) {
		err := r.DB.db.Delete(&{{$obj.Name}}{}).Error
		return err == nil, err
	}
{{end}}
`
