package main

import (
	"go-membership/app/utilis"
	"go-membership/configs"
	"go-membership/routes"
)

func main() {
	utilis.LoadEnv()
	configs.LoadConfig()
	routes.Listen()
}
