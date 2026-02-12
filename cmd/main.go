package main

import (
	"os"

	"github.com/Edge-Center/edgecentercdn-go/cmd/cli"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := cli.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
