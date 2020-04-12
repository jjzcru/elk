package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

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

	domain := "localhost"

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
	var content string

	if isQueryEnable {
		http.Handle("/playground", playground.Handler("GraphQL Playground", endpoint))
		if port == 80 {
			content = aurora.Bold(aurora.Cyan(fmt.Sprintf("http://%s/playground", domain))).String()
		} else {
			content = aurora.Bold(aurora.Cyan(fmt.Sprintf("http://%s:%d/playground", domain, port))).String()
		}

		fmt.Println(fmt.Sprintf("GraphQL playground: %s", content))
	}

	http.Handle(endpoint, srv)

	fmt.Println(strings.Join([]string{
		aurora.Bold("Server running on port").String(),
		aurora.Bold(aurora.Green(fmt.Sprintf("%d 🚀", port))).String(),
	}, " "))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
