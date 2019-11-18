package templates

var Graphql = `{{range $obj := .Model.Objects}}
  # {{.Name}}接口字段
  fragment {{$obj.PluralName}}Fields on {{.Name}} { {{range $col := $obj.Columns}}
    {{$col.Name}}{{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}
    {{$rel.Name}} {
      ...{{$rel.Target.Name}}sFields
    }{{end}}{{end}}
  }
  # 列表
  query {{.Name}}s ($currentPage: Int = 1, $perPage: Int = 20, $sort: [{{.Name}}SortType!], $search: String, $filter: {{.Name}}FilterType) {
    {{$obj.LowerName}}s(current_page: $currentPage, per_page: $perPage, sort: $sort, q: $search, filter: $filter) {
      data { {{range $col := $obj.Columns}}
          {{$col.Name}}{{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}
          {{$rel.Name}} {
            ...{{$rel.Target.Name}}sFields
        }{{end}}{{end}}
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
`

var GraphqlApi = `[
  {{range $obj := .Model.Objects}}
    {
      "title": "{{$obj.EntityName}}",
      "fields": [
        {{range $col := $obj.Columns}}
        { "name": "{{$col.Name}}", "desc": "{{$col.GetComment}}", "type": "{{$col.GetType}}", "required": "{{$col.IsRequired}}", "validator": "password", "remark": "{{$col.GetRemark}}" },
        {{end}}
      ],
    },
  {{end}}
]
`