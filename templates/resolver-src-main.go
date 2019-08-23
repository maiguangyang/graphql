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

// 自定义登录方法
func (r *MutationResolver) Login(ctx context.Context, email string) (*interface{}, error) {
	// 根据条件查询用户
	var opts gen.QueryUserHandlerOptions
	opts.Filter = &gen.UserFilterType{
		Email: &email,
	}
	user, _ := gen.QueryUserHandler(ctx, r.GeneratedResolver, opts)

	// 生成JWT Token
  ip := ctx.Value("RemoteIp")
  token := middleware.SetToken(map[string]interface{}{
    "id": user.ID,
  }, utils.EncryptMd5(ip.(string) + middleware.SecretKey["admin"].(string)), "admin")

	// 组装返回数据
	var resData interface{}
	resData = map[string]interface{}{
		"user": map[string]interface{}{
			"id"    : user.ID,
			"email" : user.Email,
			"state" : user.State,
		},
		"token": token,
	}

	return &resData, nil
}
`
