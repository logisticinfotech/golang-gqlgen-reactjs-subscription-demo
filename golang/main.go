package main

import (
	graph "golang-gqlgen-reactjs-subscription-demo/golang/app/graph"
	resolver "golang-gqlgen-reactjs-subscription-demo/golang/app/resolver"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	gqlopentracing "github.com/99designs/gqlgen/opentracing"
	"github.com/gorilla/websocket"
)

const defaultPort = "8888"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	// http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}})))
	http.Handle("/query", corsAccess(handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}),

		handler.ResolverMiddleware(gqlopentracing.ResolverMiddleware()),
		handler.RequestMiddleware(gqlopentracing.RequestMiddleware()),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}))),
	)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func corsAccess(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == "OPTIONS" {
			return
		}
		next(w, req)
	})
}
