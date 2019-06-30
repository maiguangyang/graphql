package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/novacloudcz/graphql-orm/templates"

	"github.com/inloop/goclitools"

	"gopkg.in/yaml.v2"

	"github.com/novacloudcz/graphql-orm/model"
	"github.com/urfave/cli"
)

var initCmd = cli.Command{
	Name:  "init",
	Usage: "initialize new project",
	Action: func(ctx *cli.Context) error {
		fmt.Println("Initializing project...")

		if !fileExists("graphql-orm.yml") {
			if err := createConfigFile(); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := createDummyModelFile(); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := createMainFile(); err != nil {
			return cli.NewExitError(err, 1)
		}

		if !fileExists("makefile") {
			wantCreateMakefile := goclitools.Prompt("Create makefile for run/generate commands? [y/N]")
			if strings.ToLower(wantCreateMakefile) == "y" {
				if err := createMakeFile(); err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		}

		if !fileExists("Dockerfile") {
			wantCreateDockerfile := goclitools.Prompt("Create Dockerfile for building docker images? [y/N]")
			if strings.ToLower(wantCreateDockerfile) == "y" {
				if err := createDockerFile(); err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		}

		if err := runGenerate(); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	},
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return true
	}
	return false
}

func createConfigFile() error {
	defaultPackagePath := ""
	if os.Getenv("GOPATH") != "" {
		cw, _ := os.Getwd()
		defaultPackagePath, _ = filepath.Rel(os.Getenv("GOPATH")+"/src", cw)
	}
	packagePath := goclitools.Prompt(fmt.Sprintf("Package path (default %s)", defaultPackagePath))
	if packagePath != "" {
		defaultPackagePath = packagePath
	}
	c := model.Config{Package: defaultPackagePath}

	content, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("graphql-orm.yml", content, 0644)
	return err
}
func createMainFile() error {
	c, err := model.LoadConfig()
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"%s/gen"
)

const (
	defaultPort = "80"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	urlString := os.Getenv("DATABASE_URL")
	if urlString == "" {
		panic(fmt.Errorf("missing DATABASE_URL environment variable"))
	}

	db := gen.NewDBWithString(urlString)
	defer db.Close()
	db.AutoMigrate()

	gqlHandler := handler.GraphQL(gen.NewExecutableSchema(gen.Config{Resolvers: &gen.Resolver{DB: db}}))
	playgroundHandler := handler.Playground("GraphQL playground", "/graphql")
	http.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		principalID := getPrincipalID(req)
		ctx := context.WithValue(req.Context(), gen.KeyPrincipalID, principalID)
		req = req.WithContext(ctx)
		if req.Method == "GET" {
			playgroundHandler(res, req)
		} else {
			gqlHandler(res, req)
		}
	})

	http.HandleFunc("/healthcheck", func(res http.ResponseWriter, req *http.Request) {
		if err := db.Ping(); err != nil {
			res.WriteHeader(400)
			res.Write([]byte("ERROR"))
			return
		}
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPrincipalID(req *http.Request) string {
	return req.Header.Get("principal-id")
}

`, c.Package)
	return ioutil.WriteFile("main.go", []byte(content), 0644)
}
func createDummyModelFile() error {
	content := `type User {
	email: String
	firstName: String
	lastName: String

	tasks: [Task!]! @relationship(inverse:"assignee")
}

type Task {
	title: String
	completed: Boolean
	dueDate: Time

	assignee: User @relationship(inverse:"tasks")
}
`

	if _, err := os.Stat("model.graphql"); !os.IsNotExist(err) {
		return nil
	}
	return ioutil.WriteFile("model.graphql", []byte(content), 0644)
}
func createMakeFile() error {
	content := `generate:
	go run github.com/novacloudcz/graphql-orm

run:
	DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go

voyager:
	docker run --rm -v ` + "`" + `pwd` + "`" + `/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager
`
	return ioutil.WriteFile("makefile", []byte(content), 0644)
}
func createDockerFile() error {
	c, err := model.LoadConfig()
	if err != nil {
		return err
	}
	data := TemplateData{nil, &c}
	return writeTemplate(templates.Dockerfile, "Dockerfile", data)
}

func runGenerate() error {
	return goclitools.RunInteractive("go run github.com/novacloudcz/graphql-orm")
}
