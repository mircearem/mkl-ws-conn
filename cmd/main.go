package main

import (
	"log"

	"github.com/mircearem/mklwsconn/api"
)

func main() {
	app := api.NewServer()
	log.Fatal(app.Run())
}
