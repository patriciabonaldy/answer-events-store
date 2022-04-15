package main

import (
	"log"

	"github.com/patriciabonaldy/bequest_challenge/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
