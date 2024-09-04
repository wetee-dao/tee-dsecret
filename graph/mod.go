package graph

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"wetee.app/dsecret/dkg"
	"wetee.app/dsecret/util"
)

var dkgIns *dkg.DKG

// 启动GraphQL服务器
// StartServer starts the GraphQL server.
func StartServer(d *dkg.DKG) {
	dkgIns = d
	port := util.GetEnvInt("GQL_PORT", 8880)

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

	if util.IsFileExists("./ssl/ser.pem") && util.IsFileExists("./ssl/ser.key") {
		log.Printf("connect to https://0.0.0.0:%s/ for GraphQL playground", fmt.Sprint(port))
		http.ListenAndServeTLS(":"+fmt.Sprint(port), "./ssl/ser.pem", "./ssl/ser.key", router)
	} else {
		log.Printf("connect to http://0.0.0.0:%s/ for GraphQL playground", fmt.Sprint(port))
		http.ListenAndServe(":"+fmt.Sprint(port), router)
	}
}
