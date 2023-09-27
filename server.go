package main

import (
	config "raihpeduli/configs"
	route "raihpeduli/routes"
)


func main() {
	config.ConnectDatabase()

	route.RunServer()
}