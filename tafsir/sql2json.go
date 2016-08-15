package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	BookId    string `json: "book_id"`
	BookTitle string `json: "book_title"`
	Info      string `json: "info"`
	BookInfo  string `json: "book_info"`
	Author    string `json: "author"`
	AuthorBio string `json: "author_bio"`
	Category  string `json: "category"`
	Died      string `json: "died"`
}

type Page struct {
	PageId     string    `json: "page_id"`
	PageBody   string    `json: "page_body"`
	Volume     string    `json: "volume"`
	PageNumber string    `json: "page_number"`
	Chapters   []Chapter `json: "chapters"`
	Book       Book      `json: "book"`
}

type Chapter struct {
	Heading      string `json: "heading"`
	HeadingLevel string `json: "heading_level"`
	SubLevel     string `json: "sub_level"`
	PageId       string `json: "page_id"`
}

const help = `
 Please use one of the following arguments:

 -s		Export data to JSON file in json directory that will be created
 -i		Index data into ElasticSearch

 `

var (
	saveJson = flag.Bool("s", false, "Wishing to save data to json")
	indexDB  = flag.Bool("i", false, "Indexing data to Elasticsearch")
)

func init() {
	flag.Parse()
}

func index(db *sql.DB, id string) bool {
	_, err := getPages(db, id)
	if err != nil {
		log.Println(err)
		return false
	}

	// marshall into json and append to it also the
	// bulk index meta data

	return true
}

func getBook(db *sql.DB, id string) (Book, error) {

	var (
		BookIdValue    driver.Value
		BookTitleValue driver.Value
		InfoValue      driver.Value
		BookInfoValue  driver.Value
		AuthorValue    driver.Value
		AuthorBioValue driver.Value
		CategoryValue  driver.Value
		DiedValue      driver.Value

		book Book
	)

	rows, err := db.Query("SELECT BkId, Bk, Betaka, Inf, Auth, AuthInf, cat, AD FROM main")
	if err != nil {
		log.Println(err)
		return book, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			bookid    sql.NullString
			booktitle sql.NullString
			info      sql.NullString
			bookinfo  sql.NullString
			author    sql.NullString
			authorbio sql.NullString
			category  sql.NullString
			died      sql.NullString
		)

		if err := rows.Scan(&bookid, &booktitle, &info, &bookinfo, &author, &authorbio, &category, &died); err != nil {
			log.Println(err)
			return book, err
		}

		if bookid.Valid {
			BookIdValue, err = bookid.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			BookIdValue = "null"
		}

		if booktitle.Valid {
			BookTitleValue, err = booktitle.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			BookTitleValue = "null"
		}

		if info.Valid {
			InfoValue, err = info.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			InfoValue = "null"
		}

		if bookinfo.Valid {
			BookInfoValue, err = bookinfo.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			BookInfoValue = "null"
		}

		if author.Valid {
			AuthorValue, err = author.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			AuthorValue = "null"
		}

		if authorbio.Valid {
			AuthorBioValue, err = authorbio.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			AuthorBioValue = "null"
		}

		if category.Valid {
			CategoryValue, err = category.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			CategoryValue = "null"
		}

		if died.Valid {
			DiedValue, err = died.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}
		} else {
			DiedValue = "null"
		}

		book = Book{
			BookId:    BookIdValue.(string),
			BookTitle: BookTitleValue.(string),
			Info:      InfoValue.(string),
			BookInfo:  BookInfoValue.(string),
			Author:    AuthorValue.(string),
			AuthorBio: AuthorBioValue.(string),
			Category:  CategoryValue.(string),
			Died:      DiedValue.(string),
		}

	}

	return book, nil
}

func getChapters(db *sql.DB, id string) ([]Chapter, error) {

	var (
		HeadingValue      driver.Value
		HeadingLevelValue driver.Value
		SubLevelValue     driver.Value
		PageIdValue       driver.Value

		chapters []Chapter
	)

	var trim string

	if strings.HasPrefix(id, "00") {
		trim = "00"
	} else if strings.HasPrefix(id, "0") {
		trim = "0"
	}

	id = strings.TrimPrefix(id, trim)

	rows, err := db.Query("SELECT tit, lvl, sub, id from t" + id)
	if err != nil {
		log.Println(err)
		return chapters, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			heading      sql.NullString
			headinglevel sql.NullString
			sublevel     sql.NullString
			pageid       sql.NullString
		)

		if err := rows.Scan(&heading, &headinglevel, &sublevel, &pageid); err != nil {
			log.Println(err)
			return chapters, err
		}

		if heading.Valid {
			HeadingValue, err = heading.Value()
			if err != nil {
				log.Println(err)
				return chapters, err
			}
		} else {
			HeadingValue = "null"
		}

		if headinglevel.Valid {
			HeadingLevelValue, err = headinglevel.Value()
			if err != nil {
				log.Println(err)
				return chapters, err
			}
		} else {
			HeadingLevelValue = "null"
		}

		if sublevel.Valid {
			SubLevelValue, err = sublevel.Value()
			if err != nil {
				log.Println(err)
				return chapters, err
			}
		} else {
			SubLevelValue = "null"
		}

		if pageid.Valid {
			PageIdValue, err = pageid.Value()
			if err != nil {
				log.Println(err)
				return chapters, err
			}
		} else {
			PageIdValue = "null"
		}

		chapters = append(chapters, Chapter{

			Heading:      HeadingValue.(string),
			HeadingLevel: HeadingLevelValue.(string),
			SubLevel:     SubLevelValue.(string),
			PageId:       PageIdValue.(string),
		})

	}

	return chapters, nil

}

func getPages(db *sql.DB, id string) ([]Page, error) {

	var (
		PageIdValue     driver.Value
		PageBodyValue   driver.Value
		VolumeValue     driver.Value
		PageNumberValue driver.Value

		book     Book
		pages    []Page
		chapters []Chapter

		page Page
		f    *os.File
	)

	var trim string

	if strings.HasPrefix(id, "00") {
		trim = "00"
	} else if strings.HasPrefix(id, "0") {
		trim = "0"
	}

	newid := strings.TrimPrefix(id, trim)

	rows, err := db.Query("SELECT id, nass, part, page FROM b" + newid)
	if err != nil {
		log.Println(err)
		return pages, err
	}

	defer rows.Close()

	book, err = getBook(db, id)
	if err != nil {

		log.Println(err)
		return pages, err
	}

	chapters, err = getChapters(db, id)
	if err != nil {
		log.Println(err)
		return pages, err
	}

	if *saveJson == true {
		f, err = os.Create("json/" + newid + ".json")
		if err != nil {
			return pages, err
		}

		defer f.Close()
	}

	if *indexDB == true {
		es, err := elastic.NewClient()
		if err != nil {
			log.Println(err)
			return pages, err
		}
	}

	for rows.Next() {

		var (
			pageid     sql.NullString
			pagebody   sql.NullString
			volume     sql.NullString
			pagenumber sql.NullString
		)
		if err := rows.Scan(&pageid, &pagebody, &volume, &pagenumber); err != nil {
			log.Println(err)
			return pages, err
		}

		if pageid.Valid {
			PageIdValue, err = pageid.Value()
			if err != nil {
				log.Println(err)
				return pages, err
			}
		} else {
			PageIdValue = "null"
		}

		if pagebody.Valid {
			PageBodyValue, err = pagebody.Value()
			if err != nil {
				log.Println(err)
				return pages, err
			}
		} else {
			PageBodyValue = "null"
		}

		if volume.Valid {
			VolumeValue, err = volume.Value()
			if err != nil {
				log.Println(err)
				return pages, err
			}
		} else {
			VolumeValue = "null"
		}

		if pagenumber.Valid {
			PageNumberValue, err = pagenumber.Value()
			if err != nil {
				log.Println(err)
				return pages, err
			}
		} else {
			PageNumberValue = "null"
		}

		page = Page{
			PageId:     PageIdValue.(string),
			PageBody:   PageBodyValue.(string),
			Volume:     VolumeValue.(string),
			PageNumber: PageNumberValue.(string),
			Chapters:   chapters,
			Book:       book,
		}

		pages := append(pages, page)

		jsonByte, err := json.Marshal(page)
		if err != nil {
			log.Println(err)
			return pages, err
		}

		if *saveJson == true {
			_, err = f.Write([]byte(fmt.Sprintf("%s\n", string(jsonByte))))
			if err != nil {
				log.Println(err)
				return pages, err
			}
		}

		if *indexDB == true {
			// index each page
			r, err := es.Index.BodyJson(jsonByte).Do()
			if err != nil {
				log.Println(err)
				return pages, err
			}

			resp, err := json.Marshal(r)
			if err != nil {
				log.Println(err)
				return pages, err
			}

			fmt.Println(resp)

		}
	}

	return pages, nil
}

func main() {

	if *saveJson == false && *indexDB == false {
		fmt.Println(help)
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

	if *saveJson == true {
		err = os.MkdirAll("json", 0755)
		if err != nil {
			return
		}
	}

	for i, file := range files {

		db, err := sql.Open("sqlite3", filepath.Join("db", file.Name()))
		if err != nil {
			log.Println(err)
		}

		id := strings.Split(file.Name(), ".")

		ok := index(db, id[0])
		if !ok {
			log.Println("Failed index")
		}

		db.Close()
		fmt.Printf("\n[%v] \t -  ===========> \t completed \t %s", i, file.Name())
	}

}
