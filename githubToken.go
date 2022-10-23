package main

import (
	"fmt"
	"os"
)

func githubToken() string {

	token := os.Getenv("GH_TOKEN_TRIGGER")
	if token == "" {
		panic(fmt.Errorf("missing token in GH_TOKEN_TRIGGER"))
	}
	return token
}
