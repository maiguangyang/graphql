package templates

var Graphql = `{{range $obj := .Model.Objects}}
  # {{.Name}}接口字段
  fragment {{$obj.LowerName}}sFields on {{.Name}} { {{range $col := $obj.Columns}}
    {{$col.Name}}{{end}}{{range $rel := $obj.Relationships}}{{if $rel.IsToMany}}
    {{$rel.Name}} {
      ...{{$obj.LowerName}}{{$rel.MethodName}}Fields
    }{{end}}{{end}}
  }
  # 列表
  query {{.Name}}s ($currentPage: Int = 1, $perPage: Int = 20, $sort: [{{.Name}}SortType!], $search: String, $filter: {{.Name}}FilterType) {
    {{$obj.LowerName}}s(current_page: $currentPage, per_page: $perPage, sort: $sort, q: $search, filter: $filter) {
      data {
        ...{{$obj.LowerName}}sFields
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
      ...{{$obj.LowerName}}sFields
    }
  }
  # 新增
  mutation {{.Name}}Add ($data: {{.Name}}CreateInput!) {
    create{{.Name}}(input: $data) {
      ...{{$obj.LowerName}}sFields
    }
  }

  # 修改
  mutation {{.Name}}Edit ($id: ID!, $data: {{.Name}}UpdateInput!) {
    update{{.Name}}(id: $id, input: $data) {
      ...{{$obj.LowerName}}sFields
    }
  }

  # 删除
  mutation {{.Name}}Delete ($id: ID!) {
    delete{{.Name}}(id: $id) {
      ...{{$obj.LowerName}}sFields
    }
  }

{{end}}
`
