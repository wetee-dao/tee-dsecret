package graph

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"wetee.app/dsecret/internal/dkg"
	"wetee.app/dsecret/internal/util"
)

var dkgIns *dkg.DKG

// 启动GraphQL服务器
// StartServer starts the GraphQL server.
func StartServer(d *dkg.DKG, port int) {
	dkgIns = d

	// 创建路由
	router := chi.NewRouter()

	// 添加认证中间件
	router.Use(AuthMiddleware())

	// 添加跨域中间件
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// graphql playground
	router.Handle("/", playground.Handler("WeTEE-DSECRET", "/gql"))
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers:  &Resolver{},
		Directives: NewDirectiveRoot(),
	}))

	// main graphql
	router.Handle("/gql", srv)

	if util.IsFileExists("./chain_data/ssl/ser.pem") && util.IsFileExists("./chain_data/ssl/ser.key") {
		util.LogWithRed("GraphQL playground", "https://0.0.0.0:"+fmt.Sprint(port))
		http.ListenAndServeTLS(":"+fmt.Sprint(port), "./chain_data/ssl/ser.pem", "./chain_data/ssl/ser.key", router)
	} else {
		util.LogWithRed("GraphQL playground", "http://0.0.0.0:"+fmt.Sprint(port))
		http.ListenAndServe(":"+fmt.Sprint(port), router)
	}
}
