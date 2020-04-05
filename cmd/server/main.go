package main

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/server"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	err := server.Start(port)
	if err != nil {
		fmt.Print(err.Error())
	}
}
