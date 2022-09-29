package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type httpserver struct {
}

func (s *httpserver) handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	return
}

func Newhttpserver(addr string) *http.Server {
	httpsrc := &httpserver{}
	router := mux.NewRouter()
	router.HandleFunc("/", httpsrc.handlePost)
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func main() {
	h := Newhttpserver(":8081")
	h.ListenAndServe()
}
