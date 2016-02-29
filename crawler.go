package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Crawler struct {
	// slice of strings consisting of direct urls for the books
	// later to be used for download
	books []string

	// slice of strings consisting of the urls for the pages of the individual books
	// which holds information for each book and link to pdf if it exists
	pages []string

	// slice of strings of the urls for the categories
	categories []string

	// method to be used, which at the moment consists of either scraping
	// the entire website or updating
	method string

	// which url to get. This may be updated through the code as the crawler
	// gets a different HTML document each time.
	url string
}

//ServeHTTP handling incoming requests
func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: Future functionality for inter-services communication

	// handle any incoming requests
	return
}

func (c *Crawler) init() (*os.File, error) {
	file, err := os.Create("urls.json")
	if err != nil {
		return nil, err
	}

	return file, nil
}

// run ...
func (c *Crawler) run() {
	// This code in here needs to run in its own for loop

	for {

		c.Method()

		// create new Tag
		t := new(Tag)
		t = t.New("a")

		if c.method == "scrape" {

			// Crawl through the urls of the books
			c.Crawl(c.books)
		}

		if c.method == "update" {
			s := new(Status)
			err := s.Poll()
			if err != nil {
				log.Println(err)

			}
		}

		time.Sleep(*interval)

	}
}

// Crawl starts to crawl through a given urls extracting individual book urls
func (c *Crawler) Crawl(urls []string) {

	// init file
	file, err := c.init()
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}
	ok := c.Save(file, urls)
	if !ok {
		log.Println("Could not save the urls to the file")
		return
	}
}

// Save will save each link to the top of the file
func (c *Crawler) Save(file *os.File, urls []string) (ok bool) {

	ok = true
	// TODO: Go through all pages and crawl and pull out urls of the books
	for _, url := range urls {
		// store each url into the file
		_, err := file.WriteString("http://www.shamela.ws" + url + "\n")
		if err != nil {
			ok = false
			return ok
		}
	}

	return ok
}

func (c *Crawler) Method() {

	//method of crawling
	c.method = *method

	switch c.method {

	default:
		c.url = "http://www.shamela.ws"
		break
	case "scrape":
		c.url = "http://www.shamela.ws/index.php/categories"
	}

}

func (c *Crawler) Get() (bytes []byte, err error) {

	client := new(http.Client)
	resp, err := client.Get(c.url)
	if err != nil {
		return nil, err
	}
	// is it neccessary?
	resp.Close = true

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *Crawler) Parse(t *Tag) (books []string, err error) {

	bytes, err := c.Get()
	// parse the HTML document
	re, ok := t.Compile(t.Name)
	t.Regex = re
	// TODO: test to ensure error is indeed being returned.
	if !ok {

		return nil, errors.New("could not compile the regex properly")
	}

	// Checking if there is a match
	// and if there is a match we get a []string
	c.books, err = t.Match(t.Regex, string(bytes))
	if err != nil {
		return nil, err
	}

	return c.books, nil
}
