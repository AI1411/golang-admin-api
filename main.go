package main

import (
	"api/db"
	"api/router"
)

func main() {
	dbConn := db.Init()
	router.Router(dbConn)
}
