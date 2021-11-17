package main

import (
	"fmt"
	"log"
	"qrgo/server"
)

func main() {
	srv := server.NewHTTPServer("0.0.0.0:" + "8080")
	fmt.Println("Listening on 8080")
	log.Fatal(srv.ListenAndServe())
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // lshortfile gives line number
}
