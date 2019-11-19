package templates

var Graphql = `{{range $obj := .Model.Objects}}
  # {{.Name}} 接口字段
  fragment {{$obj.PluralName}}Fields on {{.Name}} {
    {{range $col := $obj.Columns}}{{$col.Name}}
    {{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}{{$rel.Name}} {
      ...{{$rel.Target.Name}}sFields
    }
    {{end}}{{end}}
  }
  # 列表
  query {{.Name}}s ($currentPage: Int = 1, $perPage: Int = 20, $sort: [{{.Name}}SortType!], $search: String, $filter: {{.Name}}FilterType) {
    {{$obj.LowerName}}s(current_page: $currentPage, per_page: $perPage, sort: $sort, q: $search, filter: $filter) {
      data {
        {{range $col := $obj.Columns}}{{$col.Name}}
        {{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}{{$rel.Name}} {
          ...{{$rel.Target.Name}}sFields
        }
        {{end}}{{end}}
      }
      current_page
      per_page
      total
      total_page
    }
  }
  # 详情
  query {{.Name}}Detail ($id: ID, $search: String, $filter: {{.Name}}FilterType) {
    {{$obj.LowerName}}(id: $id, q: $search, filter: $filter) {
      ...{{$obj.PluralName}}Fields
    }
  }
  # 新增
  mutation {{.Name}}Add ($data: {{.Name}}CreateInput!) {
    create{{.Name}}(input: $data) {
      ...{{$obj.PluralName}}Fields
    }
  }

  # 修改
  mutation {{.Name}}Edit ($id: ID!, $data: {{.Name}}UpdateInput!) {
    update{{.Name}}(id: $id, input: $data) {
      ...{{$obj.PluralName}}Fields
    }
  }

  # 删除
  mutation {{.Name}}Delete ($id: ID!) {
    delete{{.Name}}(id: $id) {
      ...{{$obj.PluralName}}Fields
    }
  }
{{end}}
{{range $ext := .Model.ObjectExtensions}}{{$obj := $ext.Object}}
  {{range $col := $obj.Fields}}
  # {{$col.LowerName}} 接口
  {{$obj.LowerName}} {{$col.LowerName}} ({{$col.Arguments}}) {
    {{$col.Name}}({{$col.Inputs}})
  }{{end}}{{end}}
`

var GraphqlApi = `[{{range $obj := .Model.Objects}}
  {
    "title": "{{$obj.EntityName}}",
    "name": "{{$obj.LowerName}}",
    "fields": [
      {{range $col := $obj.Columns}}{ "name": "{{$col.Name}}", "desc": "{{$col.GetComment}}", "type": "{{$col.GetType}}", "required": "{{$col.IsRequired}}", "validator": "{{$col.GetValidator}}", "remark": "{{$col.GetRemark}}" },
      {{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}{ "name": "{{$rel.Name}}", "desc": "{{$rel.Target.Name}}连表查询", "type": "relationship", "required": "", "validator": "", "remark": "{{$rel.LowerName}}" },{{end}}
      {{end}}
    ],
    "data": [
      { "title": "列表", "api": "{{$obj.LowerName}}s" },
      { "title": "详情", "api": "{{$obj.LowerName}}" },
      { "title": "新增", "api": "create{{$obj.Name}}" },
      { "title": "修改", "api": "update{{$obj.Name}}" },
      { "title": "删除", "api": "delete{{$obj.Name}}" }
    ]
  },
{{end}}
{{range $ext := .Model.ObjectExtensions}}{{$obj := $ext.Object}}{{range $col := $obj.Fields}}
  {
    "title": "{{$col.Name}}",
    "name": "{{$col.Name}}",
    "fields": [
      {{$col.Fields}}    ]
  },{{end}}
{{end}}]
`