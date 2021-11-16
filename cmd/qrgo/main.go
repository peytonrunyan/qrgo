package main

import (
	"fmt"
	"log"
	"path/filepath"
	"qrgo/server"

	"github.com/joho/godotenv"
)

func main() {
	// db, err := sqlserver.InitDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	srv := server.NewHTTPServer(":" + "8082")
	fmt.Println("Listening on 8082")
	log.Fatal(srv.ListenAndServe())
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // lshortfile gives line number
	err := godotenv.Load(filepath.FromSlash(".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
