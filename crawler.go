package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Crawler struct {
	Books []string
	Pages []string
}

type Tag struct {
	Name  string
	Class []string
	Id    string
	Text  string
}

//ServeHTTP handling incoming requests
func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: Future functionality for inter-services communication

	// handle any incoming requests
	return
}

func (c *Crawler) run() {

	client := new(http.Client)
	resp, err := client.Get("http://www.shamela.ws")
	if err != nil {

		log.Println(err)
		return
	}
	resp.Close = true

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: Refactor the code so that the regex matches an entire element
	// and returns a list of elements.
	log.Println("Compiling regex")
	re := regexp.MustCompile(`"\/index.php\/book\/\d+"`)
	log.Println("Printing out the results....")
	fmt.Println(re.FindAllString(string(bytes), -1))
}

// Crawl starts to crawl through a given urls extracting individual book urls
func (c *Crawler) Crawl(urls []string) {
	// TODO: Go through all pages and crawl and pull out urls of the books
}

//New inits a new Tag (Html Element)
func (t *Tag) New(name string) *Tag {

	t.Name = name
	t.Class = []string{"class", "class2"}
	return t
}
