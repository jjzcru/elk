package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/logrusorgru/aurora"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jjzcru/elk/pkg/server/graph"
	"github.com/jjzcru/elk/pkg/server/graph/generated"
)

const defaultPort = 8080

// Start graphql server
func Start(port int, filePath string, isQueryEnable bool) error {
	if port == 0 {
		port = defaultPort
	}

	// We use env variable to set the file path
	err := os.Setenv("ELK_FILE", filePath)
	if err != nil {
		return err
	}

	graph.ServerCtx = context.Background()

	// Detect an interrupt signal and cancel all the detached tasks
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			graph.CancelDetachedTasks()
			os.Exit(1)
		}
	}()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	endpoint := fmt.Sprintf("/%s", "graphql")

	if isQueryEnable {
		http.Handle("/", playground.Handler("GraphQL playground", endpoint))
	}
	http.Handle(endpoint, srv)

	var content string

	if isQueryEnable {
		content = aurora.Bold(aurora.Cyan(fmt.Sprintf("http://localhost:%d/", port))).String()
		fmt.Println(fmt.Sprintf("GraphQL playground: %s", content))
	}

	intro := aurora.Bold("Application running on port ").String()
	content = aurora.Bold(aurora.Green(fmt.Sprintf("%d 🚀", port))).String()

	fmt.Println(intro + content)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	return nil
}
