package main

import (
	"echo/routes"
	"echo/database"
)
func init(){
	database.Connectdb("todos")
}
func main() {
	// server := echo.New()
	// routes.TodoRoute(server.Group("/todos"))
	// server.Logger.Fatal(server.Start(":8080"))
	routes.TodoRoute()
}
