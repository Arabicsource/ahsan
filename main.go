package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var start = time.Now()
var urlChan chan string

type Crawler struct{}

func main() {

	c := new(Crawler)

	urlChan = c.run()

	for {
		select {
		case url := <-urlChan:
			fmt.Println(url)

		case <-time.After(time.Second * 2):
			log.Println("Exiting")
			log.Println(time.Since(start))
			return

		}
	}
}

func (c *Crawler) run() chan string {
	var err error
	var cats []string
	urlChan = make(chan string, 10)

	for {

		if cats, err = getCategories(); err != nil {
			log.Println(err)

		}

		for i, cat := range cats {

			go func(cat string, i int) {
				// fmt.Println("http://www.shamela.ws" + url)

				books, err := c.crawlCat(cat)
				if err != nil {

					log.Println(err)
					return
				}
				for _, book := range books {
					urlChan <- fmt.Sprintf("http://www.shamela.ws%s", book)
				}
			}(cat, i)

		}

		return urlChan
	}
}

// getCategories gets all the category links from the Categories page on shamela
// http://www.shamela.ws/index.php/categories
func getCategories() (cats []string, err error) {

	resp, err := getBody("http://www.shamela.ws/index.php/categories")
	if err != nil {
		return nil, err
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\/index.php\/category\/\d+`)
	cats = re.FindAllString(string(respbody), -1)

	return cats, nil

}

// crawlCat crawls the individual category page and retrieves the urls to
// the individual books' page.
func (c *Crawler) crawlCat(cat string) (books []string, err error) {
	resp, err := getBody("http://www.shamela.ws" + cat)
	if err != nil {
		return nil, err
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`\/index.php\/book\/\d+`)
	books = re.FindAllString(string(respbody), -1)

	return books, nil
}

// getBody creates a http client and makes a Get request to the given url,
// and returns a pointer to a http.Response struct
func getBody(url string) (*http.Response, error) {
	client := new(http.Client)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
