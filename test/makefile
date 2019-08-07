generate:
	go run github.com/maiguangyang/graphql
run:
	DATABASE_URL=mysql://'root:123456@tcp(localhost:3306)/graphql?charset=utf8mb4&parseTime=True&loc=Local' go run *.go
voyager:
	docker run --rm -v `pwd`/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager
