package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	start           = time.Now()
	urlChan         chan string
	allowSQLDump    = flag.Bool("sql", false, "Dump SQL into db directory")
	allowDownload   = flag.Bool("download", false, "Download files into download directory")
	allowRARExtract = flag.Bool("unrar", false, "Extract Rar files into bok directory")
	saveJSON        = flag.Bool("save-json", false, "Wishing to save data to json")
	indexDB         = flag.Bool("index", false, "Indexing data to Elasticsearch")
	file            = flag.String("file", "", "path to file containing urls")
)

func main() {
	flag.Parse()

	if *allowDownload == true {
		fmt.Println(`Scraping http://www.shamela.ws/ for URL links to shamela books in Tafsir Category.`)

		c := new(Crawler)
		urlChan := c.run()
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

				if !contains(books, url) {

					fmt.Printf("[%d] - %v \n", count, url)

					_, err := f.WriteString(fmt.Sprintf("%v\n", url))

					if err != nil {
						panic(err)
					}

					books = append(books, fmt.Sprintf("%v", url))
					count++
				}

			case <-time.After(time.Millisecond * 5000):
				log.Println("Exiting")
				log.Println(time.Since(start))
				break Loop

			}
		}

		fmt.Println("Total books found: ", count-1)
		fmt.Println(fmt.Sprintf("\n\nStarting the downloading of the books:\n\n"))

		ct := make(chan string)
		for _, book := range books {
			go func(ct chan string, book string) {
				download(ct, book)
			}(ct, book)
		}
	Loop2:
		for {
			select {
			case str := <-ct:

				fmt.Println(str)

			case <-time.After(time.Minute * 1):
				log.Println("Exiting")
				break Loop2

			}
		}
	} else if *allowRARExtract == true {
		files, err := ioutil.ReadDir("downloads")
		if err != nil {
			log.Println(err)
		}

		for _, file := range files {
			if err = extract(file.Name()); err != nil {
				log.Println(err)
				continue
			}
		}
	} else if *allowSQLDump == true {
		files, err := ioutil.ReadDir("bok")
		if err != nil {
			log.Println(err)
		}

		for _, file := range files {
			SQLFile, err := dump(file)
			if err != nil {
				//log.Fatal(err)
				cmd := exec.Command("rm", SQLFile)
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
				log.Printf("Failed following file: %s - %v", file.Name(), err)
			}
			fmt.Printf("Completed SQL file: %v\n", SQLFile)

		}
		// Check the --save-json and --index flags
	} else if *saveJSON == true || *indexDB == true {
		c := make(chan string)

		if *saveJSON == false && *indexDB == false {
			return
		}

		err := os.Setenv("MDB_JET3_CHARSET", "cp1256")
		if err != nil {
			return
		}

		files, err := ioutil.ReadDir("db")
		if err != nil {
			log.Println(err)
			return
		}

		if *saveJSON == true {
			err = os.MkdirAll("json", 0755)
			if err != nil {
				return
			}
		}

		for _, file := range files {

			go func(file os.FileInfo, c chan string) {

				db, err := sql.Open("sqlite3", filepath.Join("db", file.Name()))
				if err != nil {
					log.Println(err)
				}

				id := strings.Split(file.Name(), ".")

				ok := index(db, id[0], c)
				if !ok {
					log.Println("Failed indexing " + file.Name())
				}

				db.Close()

			}(file, c)

		}

		for {

			select {

			case msg := <-c:
				fmt.Println(msg)

			case <-time.After(10 * time.Minute):
				return

			}

		}

	} else {
		fmt.Println(`Maktabah-cli tool v0.1-alpha-rc1 <October 24, 2016>
By aboo shayba <shaybix> aboo.shayba@gmail.com

 sql			Dump SQL into db directory")
 download		Download files into download directory")
 unrar			Extract Rar files into bok directory")
 save-json		Wishing to save data to json")
 index			Indexing data to Elasticsearch")
 file			Give file containing urls of files to download

For any support please contact the author at aboo.shayba@gmail.com		
		`)
	}

}
