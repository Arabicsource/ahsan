package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Crawler struct {
	urls []string
}

type Tag struct {
	Name  string
	Count int
	Class []string
	Id    string
	Text  string
}

//ServeHTTP handling incoming requests
func (c *Crawler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// handle any incoming requests

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
	log.Println("Compiling regex")
	re := regexp.MustCompile(`"\/index.php\/book\/\d+"`)
	log.Println("Printing out the results....")
	fmt.Println(re.FindAllString(string(bytes), -1))
}

func (t *Tag) New(name string) *Tag {

	t.Name = name
	t.Class = []string{"class", "class2"}
	return t
}
