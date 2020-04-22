package main

import (
	log "./lib/logs"
	models "./models"
	routes "./routes"
)

func main() {
	PUERTO := "5001"
	log.PrintLog("Servicio iniciado")
	models.LoadRoles()
	models.CreateDB()
	routes.LoadRouter(PUERTO)
}
