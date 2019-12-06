package templates

var Graphql = `{{range $obj := .Model.Objects}}
  # {{$obj.EntityName}} {{$obj.Name}} 接口字段
  fragment {{$obj.PluralName}}Fields on {{$obj.Name}} {
    {{range $col := $obj.Columns}}{{$col.Name}}
    {{end}}
  }

  # 列表
  query {{$obj.Name}}s ($currentPage: Int = 1, $perPage: Int = 20, $sort: [{{$obj.Name}}SortType!], $search: String, $filter: {{$obj.Name}}FilterType) {
    {{$obj.LowerName}}s(current_page: $currentPage, per_page: $perPage, sort: $sort, q: $search, filter: $filter) {
      data {
        {{range $col := $obj.Columns}}{{$col.Name}}
        {{end}}{{range $rel := $obj.Relationships}}{{$rel.Name}} {
          ...{{$rel.Target.Name}}sFields
        }
        {{end}}
      }
      current_page
      per_page
      total
      total_page
    }
  }

  # 详情
  query {{$obj.Name}}Detail ($id: ID, $search: String, $filter: {{$obj.Name}}FilterType) {
    {{$obj.LowerName}}(id: $id, q: $search, filter: $filter) {
      {{range $col := $obj.Columns}}{{$col.Name}}
      {{end}}{{range $rel := $obj.Relationships}}{{$rel.Name}} {
        ...{{$rel.Target.Name}}sFields
      }
      {{end}}
    }
  }

  # 新增
  mutation {{$obj.Name}}Add ($data: {{$obj.Name}}CreateInput!) {
    create{{$obj.Name}}(input: $data) {
      ...{{$obj.PluralName}}Fields
    }
  }

  # 修改
  mutation {{$obj.Name}}Edit ($id: ID!, $data: {{$obj.Name}}UpdateInput!) {
    update{{$obj.Name}}(id: $id, input: $data) {
      ...{{$obj.PluralName}}Fields
    }
  }

  # 删除
  mutation {{$obj.Name}}Delete ($id: ID!) {
    delete{{$obj.Name}}(id: $id) {
      ...{{$obj.PluralName}}Fields
    }
  }
{{end}}
{{range $ext := .Model.ObjectExtensions}}{{$obj := $ext.Object}}
  {{range $col := $obj.Fields}}
  # {{$col.LowerName}} 接口
  {{$obj.LowerName}} {{$col.LowerName}}{{$col.Arguments}} {
    {{$col.Name}}{{$col.Inputs}}{{if $col.IsReadonlyType}} {
      ...{{$col.TargetType}}sFields{{if $col.TargetObject.HasAnyRelationships}}
      {{range $rol := $col.TargetObject.Relationships}}
      {{$rol.Name}} {
        {{range $rel := $rol.Target.Columns}}{{$rel.Name}}
        {{end}}{{range $oRel := $rol.Target.Relationships}}{{$oRel.Name}} {
          ...{{$oRel.Target.Name}}sFields
        }
        {{end}}
      }{{end}}{{end}}
    }{{end}}
  }{{end}}{{end}}
`

var GraphqlApi = `export default [{{range $obj := .Model.Objects}}
  {
    title: '{{$obj.EntityName}}',
    name: '{{$obj.LowerName}}s',
    type: 0,
    fields: [
      {{range $col := $obj.Columns}}{ name: '{{$col.Name}}', desc: '{{$col.GetComment}}', type: '{{$col.GetType}}', required: '{{$col.IsRequired}}', validator: '{{$col.GetValidator}}', remark: '{{$col.GetRemark}}' },
      {{end}}{{range $rel := $obj.Relationships}}{ name: '{{$rel.Name}}', desc: '{{$rel.Target.Name}}连表查询', type: 'relationship', required: '', validator: '', remark: '{{$rel.LowerName}}' },
      {{end}}
    ],
    data: [
      { title: '列表', api: '{{$obj.LowerName}}s', type: 'list', method: 'query' },
      { title: '详情', api: '{{$obj.LowerName}}', type: 'detail', method: 'query' },
      { title: '新增', api: 'create{{$obj.Name}}', type: 'add', method: 'mutation' },
      { title: '修改', api: 'update{{$obj.Name}}', type: 'edit', method: 'mutation' },
      { title: '删除', api: 'delete{{$obj.Name}}', type: 'delete', method: 'mutation' },
    ],
  },
{{end}}
{{range $ext := .Model.ObjectExtensions}}{{$obj := $ext.Object}}{{range $col := $obj.Fields}}
  {
    title: '{{$col.EntityName}}',
    name: '{{$col.Name}}',
    type: 1,
    default: {{$col.GetDefault}},
    fields: [
      {{$col.Fields}}    ],
    data: [
      { title: '详情', api: '{{$col.Name}}', type: 'detail', method: '{{$obj.LowerName}}' },
    ],
  },{{end}}
{{end}}];
`