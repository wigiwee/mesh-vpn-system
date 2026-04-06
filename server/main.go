package main

import (
	"log"
	"net/http"
	"server/config"
	"server/routers"
	"strconv"
)

func main() {

	r := routers.Router()
	log.Println("[INFO] starting the server")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), r))
	log.Println("listening at ", strconv.Itoa(config.Port))
}
