package server

import (
	"fmt"
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

	endpoint := fmt.Sprintf("/%s", "graphql")

	http.Handle("/", playground.Handler("GraphQL playground", endpoint))
	http.Handle(endpoint, srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}
