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

{{range $valid := $object.Validators}}
	{{$valid.MethodName}} {{$valid.GoType}} ` + "`" + `{{$valid.ModelTags}}` + "`" + `{{end}}

{{range $rel := $object.Relationships}}
	{{$rel.MethodName}} {{$rel.GoType}} ` + "`" + `{{$rel.ModelTags}}` + "`" + `{{end}}
}

{{end}}
`
