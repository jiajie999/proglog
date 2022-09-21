package server

import (
	"github.com/jiajie999/proglog/LetsGo/internal/server"
	"log"
)

package main
import (
"log"
"github.com/jiajie999/proglog/internal/server"
)
func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}