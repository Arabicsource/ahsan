package main

import (
	"errors"
	"fmt"
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

	// url channel
	urls chan []string
}

//ServeHTTP handling incoming requests
func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: Future functionality for inter-services communication

	// handle any incoming requests
	return
}

func (c *Crawler) init() (*os.File, error) {

	// create urls.json file
	file, err := os.Create("urls.json")
	if err != nil {
		return nil, err
	}

	return file, nil
}

// run ...
func (c *Crawler) run() {
	// This code in here needs to run in its own for loop

	var err error
	c.urls = make(chan []string, 4)

	for {

		// set the method of crawling
		// and the starting point (url)
		c.Method()

		// create new Tag
		t := new(Tag)
		t = t.New("a")

		err = c.Parse(t)
		if err != nil {
			continue
		}
		if c.method == "scrape" {

			// Crawl through the urls of the categories
			ok := c.Crawl(c.categories)
			if !ok {
				log.Println("Could not crawl through the urls of the categories")
			}
		}

		if c.method == "update" {
			s := new(Status)
			err := s.Poll()
			if err != nil {
				log.Println(err)

			}
		}

		time.Sleep(time.Second * *interval)

	}
}

// Crawl starts to crawl through a given urls extracting
// individual book urls
func (c *Crawler) Crawl(urls []string) bool {

	ok := c.Save(urls)
	if !ok {
		log.Println("Could not save the urls to the file")
		return ok
	}

	for _, url := range urls {
		go func(url string) {
			// fmt.Println("http://www.shamela.ws" + url)

			c.crawlPage(url)
		}(url)

	}
	select {
	case <-c.urls:
		fmt.Println(<-c.urls)
	}
	return ok
}

// Save will save each link to the top of the file
func (c *Crawler) Save(urls []string) (ok bool) {

	// init file
	file, err := c.init()
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}

	ok = true
	// TODO: Go through all pages and crawl and
	// pull out urls of the books
	for _, url := range urls {
		// store each url into the file
		_, err := file.WriteString("http://www.shamela.ws" + url + "\n")
		if err != nil {
			ok = false
			return ok
		}
		// log.Println("http://www.shamela.ws" + url + "\n")
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

func (c *Crawler) Get(url string) (bytes []byte, err error) {

	client := new(http.Client)
	resp, err := client.Get(url)
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

func (c *Crawler) Parse(t *Tag) (err error) {

	bytes, err := c.Get(c.url)
	// parse the HTML document
	re, ok := t.Compile(t.Name, `\/index.php\/category\/\d+`)
	t.Regex = *re

	// TODO: test to ensure error is indeed being returned.
	if !ok {

		return errors.New("could not compile the regex properly")
	}

	// Checking if there is a match
	// and if there is a match we get a []string
	c.categories, err = t.Match(t.Regex, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (c *Crawler) crawlPage(url string) {

	// Make get requests to the category page and send the links
	// through the urls channel
	rsp, err := c.Get("http://www.shamela.ws" + url)
	if err != nil {
		log.Println(err)
	}

	t := new(Tag)
	t.New("a")

	re, ok := t.Compile(t.Name, `\/index.php\/book\/\d+`)

	// TODO: test to ensure error is indeed being returned.
	if !ok {

		log.Printf("Could not compile the regex properly")
	}

	t.Regex = *re
	c.books, err = t.Match(t.Regex, string(rsp))
	if err != nil {

		log.Printf("could not match the regex to the body: %v", err)

	}
	// log.Println(c.books)

	c.urls <- c.books

}
