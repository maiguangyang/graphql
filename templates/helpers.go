package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/novacloudcz/goclitools"
	"github.com/maiguangyang/graphql/model"
)

type TemplateData struct {
	Model     *model.Model
	Config    *model.Config
	RawSchema *string
}

func WriteTemplate(t, filename string, data TemplateData) error {
	return WriteTemplateRaw(t, filename, data)
}

func WriteTemplateRaw(t, filename string, data interface{}) error {
	temp, err := template.New(filename).Parse(t)
	if err != nil {
		return err
	}
	var content bytes.Buffer
	writer := io.Writer(&content)

	err = temp.Execute(writer, &data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, content.Bytes(), 0777)
	if err != nil {
		return err
	}
	if path.Ext(filename) == ".go" {
		return goclitools.RunInteractive(fmt.Sprintf("goimports -w %s", filename))
	}
	return nil
}

// 生成前端接口文档
func WriteInterfaceTemplate(t, filename string, data TemplateData) error {
	return WriteInterfaceTemplateRaw(t, filename, data)
}

func WriteInterfaceTemplateRaw(t, filename string, data interface{}) error {
  temp, err := template.New(filename).Parse(t)
  if err != nil {
    return err
  }

  // type Inventory struct {
  //   Material string
  //   Count    uint
  // }

  // sweaters := Inventory{"wool", 17}
  // temp, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
  // if err != nil {
  //   return err
  // }

  var content bytes.Buffer
  writer := io.Writer(&content)

  err = temp.Execute(writer, &data)
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(filename, content.Bytes(), 0777)
  if err != nil {
    return err
  }
  return nil
}