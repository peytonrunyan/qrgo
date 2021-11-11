package main

import (
	"log"
	"qrgo/pkg/models/sqlserver"
	"qrgo/server"

	"github.com/joho/godotenv"
)

func main() {
	db, err := sqlserver.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := server.NewHTTPServer(":" + "8082")
	log.Fatal(srv.ListenAndServe())
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // lshortfile gives line number
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
