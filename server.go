package main

import (
	config "raihpeduli/configs"
	route "raihpeduli/routes"
)


func main() {
	config.InitDB()

	route.RunServer()
}