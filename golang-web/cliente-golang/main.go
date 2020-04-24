package main

import (
	log "./lib/logs"
	routes "./routes"
)

func main() {
	log.PrintLog("Servicio iniciado")
	routes.LoadRoles()
	routes.LoadRouter()
}
