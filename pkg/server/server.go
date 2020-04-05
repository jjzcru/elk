package server

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jjzcru/elk/pkg/server/graph"
	"github.com/jjzcru/elk/pkg/server/graph/generated"
)

const defaultPort = "8080"

func Start(port string) error {
	if port == "" {
		port = defaultPort
	}

	// We use env variable to set the file path
	/*err := os.Setenv("ELK_FILE", "")
	if err != nil {
		return err
	}*/

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground 2", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}
