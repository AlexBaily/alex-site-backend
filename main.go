package main

import (
	"log"
	"net/http"
	"strings"
)

//Http handler for responding to http/s requests.
func roothandler(w http.ResponseWriter, r *http.Request) {
	//Set response headers.
	w.Header().Add("statusDescription", "200 OK")
	w.Header().Set("statusDescription", "200 OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	URISegments := strings.Split(r.URL.Path, "/")
	w.Write([]byte(URISegments[1]))
}

func main() {
	//Create a new mux router.
	mux := http.NewServeMux()
	mux.HandleFunc("/", roothandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
