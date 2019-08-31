package templates

var ResolverSrc = `package src

import (
	"context"
	"github.com/maiguangyang/graphql/events"
	"{{.Config.Package}}/gen"
	"{{.Config.Package}}/utils"
	"{{.Config.Package}}/middleware"
)


func New(db *gen.DB, ec *events.EventController) *Resolver {
	resolver := NewResolver(db, ec)

	// resolver.Handlers.CreateUser = func(ctx context.Context, r *gen.GeneratedMutationResolver, input map[string]interface{}) (item *gen.Company, err error) {
	// 	return gen.CreateUserHandler(ctx, r, input)
	// }

	return resolver
}

// You can extend QueryResolver for adding custom fields in schema
// func (r *QueryResolver) Hello(ctx context.Context) (string, error) {
// 	return "world", nil
// }
`
