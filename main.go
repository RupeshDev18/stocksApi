package main

import (
	"fmt"
	"github.com/rupeshdev18/stocksApi/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
