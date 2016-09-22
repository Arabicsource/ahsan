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

	"github.com/ttacon/chalk"
)

// Crawler ...
type Crawler struct{}

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
	paginationURL := regexp.MustCompile(`\/index.php\/category\/\d+\/page-\d`)
	pagination := paginationURL.FindAllString(string(respbody), -1)
	if len(pagination) == 0 {
		return books, nil
	}

	maxPages, err := getLastPage(pagination[len(pagination)-1])
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Max number of pages found for this category are: ", maxPages)

	// The default category page is the first page for the category page
	// so long the number is less or equals to the maxPages for that category
	// execute a goroutine, and process those pages concurrently.
	for i := 1; i <= maxPages; i++ {

		fmt.Printf("Crawling through page number %d\n", i)

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

func download(count chan string, url string) {
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

	r, err := http.Get(book[len(book)-1])

	if err != nil {
		log.Println("could not download the book, with err : ", err)
		return
	}

	n, err := io.Copy(f, r.Body)
	if err != nil {
		log.Println("Could not copy the content to the newly created file, with err: ", err)
		return
	}

	r.Body.Close()

	count <- fmt.Sprintf(fmt.Sprintf("Downloading %v \t ....", book[len(book)-1]) + chalk.Bold.TextStyle(fmt.Sprintf("\t%v kb  Downloaded. Done!\n", n/int64(1000))))
}

func contains(urlSlice []string, val string) bool {
	for _, url := range urlSlice {
		if url == val {
			return true
		}
	}
	return false
}
