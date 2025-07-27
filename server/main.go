package main

import (
	"log"
	"net/http"
	"server/routers"
)

func main() {

	r := routers.Router()
	log.Println("server started")
	log.Fatal(http.ListenAndServe(":4000", r))
	log.Println("listening at 4000")
}
