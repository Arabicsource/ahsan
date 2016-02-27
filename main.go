package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var interval = flag.Float64("interval", 12, "Default 12 Hours")

type Crawler struct{}

// Book to be downloaded from the shamela website
type Book struct {
	Bid  int
	Name string
	Link string
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

	c := new(Crawler)

	file, err := os.Create("urls.json")

	if err != nil {
		log.Println(err)

	}

	go c.run(file)
	fmt.Println("Listening on local port 8000")
	err = http.ListenAndServe(":8000", c)
	if err != nil {
		log.Fatalln(err)

	}

}

func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// handle any incoming requests

}

func (c *Crawler) run(file *os.File) {

	fileinfo  := file.Stat()
	fileinfo.
	// a for loop
}
