package main

import (
	routes "./routes"
	log "./lib/logs"
)

func main() {
	log.PrintLog("Servicio iniciado")
	routes.LoadRouter()
}
