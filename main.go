package main

import (
	"go-project/infra"
	"go-project/setup"
)

var (
	boot infra.Bootstrap
)

func main() {
	env := boot.LoadEnv()

	boot.SetupDatabase(env)

	server := boot.RunServer()

	setup.PrepareRoutes(server)

	server.Run(":8080")
}
