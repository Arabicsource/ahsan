package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var start = time.Now()
var urlChan chan string

type Crawler struct{}

func main() {

	fmt.Println(`

Scraping http://www.shamela.ws/:

	`)

	c := new(Crawler)

	urlChan = c.run()
	defer close(urlChan)
	count := 1

	f, err := os.Create("urls.txt")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	var books []string
Loop:
	for {
		select {
		case url := <-urlChan:
			fmt.Printf("[%d] - %v \n", count, url)
			_, err := f.WriteString(fmt.Sprintf("%v\n", url))
			if err != nil {
				panic(err)
			}
			books = append(books, fmt.Sprintf("%v", url))
			count++

		case <-time.After(time.Millisecond * 5000):
			log.Println("Exiting")
			log.Println(time.Since(start))
			break Loop

		}
	}
	fmt.Println("Total books found: ", count-1)
	fmt.Println(fmt.Sprintf("\n\nStarting the downloading of the books:\n\n\n"))

	for _, book := range books {
		download(book)
	}

}

// run is the starting point
func (c *Crawler) run() chan string {
	urlChan = make(chan string)

	for {

		// get tafsir category
		cats := []string{"/index.php/category/127"}

		// loop through each category page, and launch a goroutine for each
		for i, cat := range cats {

			go func(cat string, i int, urlChan chan string) {
				// fmt.Println("http://www.shamela.ws" + url)

				// get slice of urls of books (links to pages of individual books)
				books, err := c.crawlCat(cat, urlChan)
				if err != nil {

					log.Println(err)
					return
				}
				for _, book := range books {
					urlChan <- fmt.Sprintf("http://www.shamela.ws%s", book)
				}
			}(cat, i, urlChan)

		}

		return urlChan
	}
}

// crawlCat crawls the individual category page and retrieves the urls of books  to
// the individual books' page.
func (c *Crawler) crawlCat(cat string, urlChan chan string) (books []string, err error) {
	resp, err := getBody("http://www.shamela.ws" + cat)
	if err != nil {
		return nil, err
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// regex for book url
	re := regexp.MustCompile(`\/index.php\/book\/\d+`)
	books = re.FindAllString(string(respbody), -1)

	// regex for last page url
	paginationUrl := regexp.MustCompile(`\/index.php\/category\/\d+\/page-\d`)
	pagination := paginationUrl.FindAllString(string(respbody), -1)

	if len(pagination) == 0 {
		return books, nil
	}
	maxPages, err := getLastPage(pagination[len(pagination)-1])
	if err != nil {
		log.Println(err)
	}

	// The default category page is the first page for the category page
	// so long the number is less or equals to the maxPages for that category
	// execute a goroutine, and process those pages concurrently.
	for i := 1; i <= maxPages; i++ {

		// goroutine scraping the page number n for a particular category
		go getCatPage(i, cat, urlChan)

	}

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

// getLastPage returns the last page number if it can find it,
// or returns an error.
func getLastPage(url string) (int, error) {
	last := strings.Split(url, "-")
	n := last[len(last)-1]
	return strconv.Atoi(n)
}

func getCatPage(i int, cat string, urlChan chan string) {

	resp, err := getBody("http://www.shamela.ws" + cat + "/page-" + strconv.Itoa(i))
	if err != nil {
		log.Println(err)
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	re := regexp.MustCompile(`\/index.php\/book\/\d+`)
	pBooks := re.FindAllString(string(respbody), -1)

	for _, book := range pBooks {

		urlChan <- fmt.Sprintf("http://www.shamela.ws%s", book)
	}
}

func download(url string) {

	resp, err := getBody(url)
	if err != nil {
		log.Println(err)
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()

	re := regexp.MustCompile(`http://shamela.ws/books/\d+/\d+.rar`)
	book := re.FindAllString(string(respbody), -1)
	if book == nil {
		fmt.Println("no match!")
		return
	}

	// get the filename so to replicate it locally when we create it
	bookName := strings.SplitAfter(string(book[len(book)-1]), "/")
	fileName := bookName[len(bookName)-1]

	// now we have the link to the rar file (book) and need to download it.
	// create downloads directory if it doesn't exist
	if _, err = os.Stat("downloads"); os.IsNotExist(err) {
		err = os.Mkdir("downloads", 0700)
		if err != nil {
			log.Println("could not create the directory downloads, with err: ", err)
			return
		}
	}

	var f *os.File

	// create the file if it does not exist
	if _, err = os.Stat("downloads/" + string(fileName)); os.IsNotExist(err) {
		f, err = os.Create("downloads/" + fileName)
		if err != nil {
			log.Println("could not create file, with err: ", err)
			return
		}
		defer f.Close()
	}

	fmt.Printf("Downloading %v now ....", book[len(book)-1])
	r, err := http.Get(book[len(book)-1])
	if err != nil {
		log.Println("could not download the book, with err : ", err)
		return
	}

	fmt.Printf("...Done!\n")

	defer r.Body.Close()

	n, err := io.Copy(f, r.Body)
	if err != nil {
		log.Println("Could not copy the content to the newly created file, with err: ", err)
		return
	}

	fmt.Printf("Downloaded %v number of bytes......\n", n)

}
