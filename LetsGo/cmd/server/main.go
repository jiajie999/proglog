package main

import (
	"github.com/jiajie999/proglog/LetsGo/internal/server"
	"log"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
