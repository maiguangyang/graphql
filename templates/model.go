package templates

var Model = `package gen

import (
	"time"
	"github.com/maiguangyang/graphql/resolvers"
)


{{range $object := .Model.Objects}}

type {{.Name}}ResultType struct {
	resolvers.EntityResultType
}

type {{.Name}} struct {
{{range $col := $object.Columns}}
	{{$col.MethodName}} {{$col.GoType}} ` + "`" + `{{$col.ModelTags}}` + "`" + `{{end}}

{{range $rel := $object.Relationships}}
{{$rel.MethodName}} {{$rel.GoType}} ` + "`" + `{{$rel.ModelTags}}` + "`" + `
{{end}}
}

type {{.Name}}Changes struct {
{{range $col := $object.Columns}}
	{{$col.MethodName}} {{$col.GoType}}{{end}}
	{{range $rel := $object.Relationships}}{{if $rel.IsToMany}}
	{{$rel.ChangesName}} {{$rel.ChangesType}}{{end}}{{end}}
}

{{end}}

// used to convert map[string]interface{} to EntityChanges struct
func ApplyChanges(changes map[string]interface{}, to interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		TagName:     "json",
		Result:      to,
		ZeroFields:  true,
		// This is needed to get mapstructure to call the gqlgen unmarshaler func for custom scalars (eg Date)
		DecodeHook: func(a reflect.Type, b reflect.Type, v interface{}) (interface{}, error) {

			if b == reflect.TypeOf(time.Time{}) {
				switch a.Kind() {
				case reflect.String:
					return time.Parse(time.RFC3339, v.(string))
				case reflect.Float64:
					return time.Unix(0, int64(v.(float64))*int64(time.Millisecond)), nil
				case reflect.Int64:
					return time.Unix(0, v.(int64)*int64(time.Millisecond)), nil
				default:
					return v, fmt.Errorf("Unable to parse date from %v", v)
				}
			}

			if reflect.PtrTo(b).Implements(reflect.TypeOf((*graphql.Unmarshaler)(nil)).Elem()) {
				resultType := reflect.New(b)
				result := resultType.MethodByName("UnmarshalGQL").Call([]reflect.Value{reflect.ValueOf(v)})
				err, _ := result[0].Interface().(error)
				return resultType.Elem().Interface(), err
			}

			return v, nil
		},
	})

	if err != nil {
		return err
	}

	return dec.Decode(changes)
}
`
