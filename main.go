package main

import (
	"fmt"
	"log"
	"net/http"
)

// Book to be downloaded from the shamela website
type Book struct {
	Bid  int
	Name string
	Page string
	Bok  string
	Pdf  string
}

// Link of the books that are process and stored in json file
// to be consumed by the crawler
type Link struct {

	// The address of the .bok files only as they are the ones to be downloaded
	// and processed in preparation for indexing it into Elasticsearch.
	Address string `json: "address"`
}

// Basic http server listening on port 8000
func main() {
	fmt.Println("Listening on local port 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)

	}

}
