package templates

var HTTPHandler = `package gen

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/99designs/gqlgen/handler"
	jwtgo "github.com/dgrijalva/jwt-go"
	"{{.Config.Package}}/utils"
	"{{.Config.Package}}/middleware"
	"{{.Config.Package}}/cache"
)

var RidesCache *cache.Cache
var redisErr error

func GetHTTPServeMux(r ResolverRoot, db *DB) *mux.Router {
	// mux := http.NewServeMux()
	mux := mux.NewRouter()
	mux.Use(middleware.AuthHandler)

  RidesCache, redisErr = cache.NewCache("localhost:6379", "", 24*time.Hour)
  if redisErr != nil {
    log.Fatalf("cannot create APQ redis cache: %v", redisErr)
  }

	executableSchema := NewExecutableSchema(Config{Resolvers: r})
	gqlHandler := handler.GraphQL(executableSchema,
		// 中间件进行登录Token校验
		utils.RouterIsAuthMiddleware,

    // redis缓存
    handler.EnablePersistedQueryCache(RidesCache),
	)

	loaders := GetLoaders(db)

	playgroundHandler := handler.Playground("GraphQL playground", "/graphql")
	mux.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		principalID := getPrincipalID(req)
		ctx := context.WithValue(req.Context(), KeyPrincipalID, principalID)
		ctx = context.WithValue(ctx, KeyLoaders, loaders)
		ctx = context.WithValue(ctx, KeyExecutableSchema, executableSchema)
		req = req.WithContext(ctx)
		if req.Method == "GET" {
			playgroundHandler(res, req)
		} else {
			gqlHandler(res, req)
		}
	})
	handler := mux

	return handler
}

func getPrincipalIDFromContext(ctx context.Context) *string {
	v, _ := ctx.Value(KeyPrincipalID).(*string)
	return v
}
func getJWTClaimsFromContext(ctx context.Context) *JWTClaims {
	v, _ := ctx.Value(KeyJWTClaims).(*JWTClaims)
	return v
}

func getPrincipalID(req *http.Request) *string {
	pID := req.Header.Get("principal-id")
	if pID != "" {
		return &pID
	}
	c, _ := getJWTClaims(req)
	if c == nil {
		return nil
	}
	return &c.Subject
}

type JWTClaims struct {
	jwtgo.StandardClaims
	Scope *string
}

func getJWTClaims(req *http.Request) (*JWTClaims, error) {
	var p *JWTClaims

	tokenStr := strings.Replace(req.Header.Get("authorization"), "Bearer ", "", 1)
	if tokenStr == "" {
		return p, nil
	}

	p = &JWTClaims{}
	jwtgo.ParseWithClaims(tokenStr, p, nil)
	return p, nil
}

`
