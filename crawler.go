package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Crawler struct {
	books      []string
	pages      []string
	categories []string
}

// Tag represents the HTML element that will be parsed and pulled from
// html document
type Tag struct {
	Name  string
	Class []string
	Id    string
	Text  string
	Regex *regexp.Regexp
}

//ServeHTTP handling incoming requests
func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: Future functionality for inter-services communication

	// handle any incoming requests
	return
}

func (c *Crawler) run() {
	// This code in here needs to run in its own for loop

	client := new(http.Client)
	resp, err := client.Get("http://www.shamela.ws")
	if err != nil {

		log.Println(err)
		return
	}
	// is it neccessary?
	resp.Close = true

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	//	log.Println("Compiling regex")
	//	re := regexp.MustCompile(`\/index.php\/book\/\d+`)

	// create new Tag
	t := new(Tag)
	t = t.New("a")

	re, ok := t.Compile(t.Name)
	t.Regex = re
	// TODO: test to ensure error is indeed being returned.
	if !ok {
		log.Println("Could not compile the regex properly")
	}
	// Checking if there is a match
	// and if there is a match we get a []string
	c.books, err = t.Match(t.Regex, string(bytes))
	if err != nil {
		log.Println(err)
		// time.Sleep(interval)
	}
	c.Crawl(c.books)
	// time.Sleep(interval)

}

// Crawl starts to crawl through a given urls extracting individual book urls
func (c *Crawler) Crawl(urls []string) {
	// TODO: Go through all pages and crawl and pull out urls of the books
	for _, url := range urls {
		fmt.Println("http://www.shamela.ws" + url)
	}
}

// Save will save each link to the top of the file
func (c *Crawler) Save(file *os.File, url string) (ok bool, err error) {

	// save the link to the top of the file and return ok.
	return true, nil
}

//New inits a new Tag (Html Element)
func (t *Tag) New(name string) *Tag {

	t.Name = name
	return t
}
func (t *Tag) Compile(name string) (*regexp.Regexp, bool) {

	var ok bool
	var re *regexp.Regexp
	// Checking the given element against pre-defined cases and executing relevant regex
	// to pull out that element from the HTML
	// TODO:: refactor the regex bit to its own method
	switch t.Name {
	default:
		ok = false
		re = nil

		return re, ok
	case "a":
		log.Println("Instantiating for anchor element!")
		re = regexp.MustCompile(`<a\s\.*>\.*<\/a>`)
		ok = true
		return re, ok
	}

}

// Match will match up a regular expression with a body of string (html document)
func (t *Tag) Match(re *regexp.Regexp, body string) (result []string, err error) {

	// TODO: Refactor the code so that the regex matches an entire element
	// and returns a list of elements.

	// Storing the list of books' urls in Crawler.Books
	result = re.FindAllString(string(body), -1)

	if result == nil {
		err := errors.New("Found no Match!")
		return nil, err
	}
	return result, nil
}
