package main

import (
	"os"

	"github.com/beakeyz/dadjoke-gen/pkg/server"
)

func main() {
  os.Exit(int(server.RunServer()))
}
