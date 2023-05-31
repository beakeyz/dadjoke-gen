package main

import (
	"os"

	"github.com/beakeyz/dadjoke-gen/pkg/build"
)

func main() {
	os.Exit(build.RunBuild())
}
