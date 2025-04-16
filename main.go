package main

import (
	_ "amg-backend/docs"
	"amg-backend/server"
)

// @title amg-backend
// @version 1.0
// @description AMG - AnhMy Global Kindergarten
// @in header
func main() {
	server.RegisterServer()
}
