package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/novacloudcz/goclitools"

	"github.com/maiguangyang/graphql/model"
	"github.com/maiguangyang/graphql/templates"
	"github.com/urfave/cli"
)

var genCmd = cli.Command{
	Name:  "generate",
	Usage: "generate contents",
	Action: func(ctx *cli.Context) error {
		if err := generate("model.graphql"); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

func generate(filename string) error {
	fmt.Println("Generating contents from", filename, "...")
	modelSource, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	m, err := model.Parse(string(modelSource))
	if err != nil {
		return err
	}

	c, err := model.LoadConfig()
	if err != nil {
		return err
	}

	if _, err := os.Stat("./gen"); os.IsNotExist(err) {
		os.Mkdir("./gen", 0777)
	}

	err = model.EnrichModelObjects(&m)
	if err != nil {
		return err
	}

	err = generateFiles(&m, &c)
	if err != nil {
		return err
	}

	err = model.EnrichModel(&m)
	if err != nil {
		return err
	}

	schema, err := model.PrintSchema(m)
	if err != nil {
		return err
	}

	schema = "# This schema is generated, please don't update it manually\n\n" + schema

	if err := ioutil.WriteFile("gen/schema.graphql", []byte(schema), 0644); err != nil {
		return err
	}

	fmt.Println("Running gqlgen generator...")
	if err := goclitools.RunInteractiveInDir("go run github.com/99designs/gqlgen", "./gen"); err != nil {
		return err
	}

	// for _, obj := range plainModel.Objects() {
	// 	s1 := fmt.Sprintf("type %s struct {", obj.Name())
	// 	s2 := fmt.Sprintf("type %s struct {\n\t%sExtensions", obj.Name(), obj.Name())
	// 	if err := replaceStringInFile("gen/models_gen.go", s1, s2); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func generateFiles(m *model.Model, c *model.Config) error {
	data := templates.TemplateData{Model: m, Config: c}
	if err := templates.WriteTemplate(templates.Database, "gen/database.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.GeneratedResolver, "gen/resolver.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.GQLGen, "gen/gqlgen.yml", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Model, "gen/models.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Filters, "gen/filters.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.QueryFilters, "gen/query-filters.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Keys, "gen/keys.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Callback, "gen/callback.go", data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Loaders, "gen/loaders.go", data); err != nil {
		return err
	}
	return nil
}

