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
func Start(port int, filePath string, isQueryEnable bool, token string) error {
	if port == 0 {
		port = defaultPort
	}

	domain := "localhost"

	isCorsEnable := true

	addContext := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isCorsEnable {
				allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
				if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
					w.Header().Set("Access-Control-Expose-Headers", "Authorization")
				}
			}

			ctx := context.WithValue(r.Context(), graph.ElkFileKey, filePath)
			ctx = context.WithValue(ctx, graph.TokenKey, token)
			ctx = context.WithValue(ctx, graph.AuthorizationKey, r.Header.Get("Authorization"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	graph.ServerCtx = context.Background()

	// Detect an interrupt signal and cancel all the detached tasks
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		/*select {
		case <-c:
			graph.CancelDetachedTasks()
			os.Exit(1)
		}*/

		<-c
		graph.CancelDetachedTasks()
		os.Exit(1)
	}()

	srv := addContext(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})))

	endpoint := fmt.Sprintf("/%s", "graphql")
	var content string

	if isQueryEnable {
		http.Handle("/playground", playground.Handler("GraphQL Playground", endpoint))
		if port == 80 {
			content = aurora.Bold(aurora.Cyan(fmt.Sprintf("http://%s/playground", domain))).String()
		} else {
			content = aurora.Bold(aurora.Cyan(fmt.Sprintf("http://%s:%d/playground", domain, port))).String()
		}

		fmt.Printf("GraphQL playground: %s \n", content)
	}

	http.Handle(endpoint, srv)

	if len(token) > 0 {
		fmt.Println(strings.Join([]string{
			aurora.Bold("Authorization token:").String(),
			aurora.Bold(aurora.Cyan(token)).String(),
		}, " "))
	}

	fmt.Println(strings.Join([]string{
		aurora.Bold("Server running on port").String(),
		aurora.Bold(aurora.Green(fmt.Sprintf("%d 🚀", port))).String(),
	}, " "))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
