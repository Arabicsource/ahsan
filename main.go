package main

import (
	"fmt"
	"log"
	"net/http"
)

// Basic http server listening on port 8000

func main() {
	fmt.Println("Listening on local port 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)

	}

}
