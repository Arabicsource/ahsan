package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	_ "github.com/mattn/go-sqlite3"
)

// Book ...
type Book struct {
	BookID    string `json:"book_id"`
	BookTitle string `json:"book_title"`
	Info      string `json:"info"`
	BookInfo  string `json:"book_info"`
	Author    string `json:"author"`
	AuthorBio string `json:"author_bio"`
	Category  string `json:"category"`
	Died      string `json:"died"`
}

// Page ...
type Page struct {
	PageID     string  `json:"page_id"`
	PageBody   string  `json:"page_body"`
	Volume     string  `json:"volume"`
	PageNumber string  `json:"page_number"`
	Chapter    Chapter `json:"chapter"`
	Book       Book    `json:"book"`
}

// Chapter ...
type Chapter struct {
	Heading      string `json:"heading"`
	HeadingLevel string `json:"heading_level"`
	SubLevel     string `json:"sub_level"`
	PageID       string `json:"page_id"`
}

func index(db *sql.DB, id string, c chan string) bool {
	_, err := getPages(db, id)
	if err != nil {
		log.Println(err)
		return false
	}

	// marshall into json and append to it also the
	// bulk index meta data

	c <- fmt.Sprintf("\n[%s] \t -  ===========> \t completed\n", id)
	return true
}

func getBook(db *sql.DB, id string) (Book, error) {

	var (
		BookIDValue    driver.Value
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
			BookIDValue, err = bookid.Value()
			if err != nil {
				log.Println(err)
				return book, err
			}

		} else {
			BookIDValue = "null"
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
			BookID:    BookIDValue.(string),
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

func getChapter(db *sql.DB, id, pageid string) (Chapter, error) {

	var (
		HeadingValue      driver.Value
		HeadingLevelValue driver.Value
		SubLevelValue     driver.Value
		PageIDValue       driver.Value

		chapter Chapter
	)

	var trim string

	if strings.HasPrefix(id, "00") {
		trim = "00"
	} else if strings.HasPrefix(id, "0") {
		trim = "0"
	}

	id = strings.TrimPrefix(id, trim)

	rows, err := db.Query("SELECT tit, lvl, sub, id FROM t" + id + " WHERE id <= '" + pageid + "' ORDER BY id DESC LIMIT 1")
	if err != nil {
		log.Println(err)
		return chapter, err
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
			return chapter, err
		}

		if heading.Valid {
			HeadingValue, err = heading.Value()
			if err != nil {
				log.Println(err)
				return chapter, err
			}
		} else {
			HeadingValue = "null"
		}

		if headinglevel.Valid {
			HeadingLevelValue, err = headinglevel.Value()
			if err != nil {
				log.Println(err)
				return chapter, err
			}
		} else {
			HeadingLevelValue = "null"
		}

		if sublevel.Valid {
			SubLevelValue, err = sublevel.Value()
			if err != nil {
				log.Println(err)
				return chapter, err
			}
		} else {
			SubLevelValue = "null"
		}

		if pageid.Valid {
			PageIDValue, err = pageid.Value()
			if err != nil {
				log.Println(err)
				return chapter, err
			}
		} else {
			PageIDValue = "null"
		}

		chapter = Chapter{

			Heading:      HeadingValue.(string),
			HeadingLevel: HeadingLevelValue.(string),
			SubLevel:     SubLevelValue.(string),
			PageID:       PageIDValue.(string),
		}

	}

	return chapter, nil

}

func getPages(db *sql.DB, id string) ([]Page, error) {

	var (
		PageIDValue     driver.Value
		PageBodyValue   driver.Value
		VolumeValue     driver.Value
		PageNumberValue driver.Value

		book    Book
		pages   []Page
		chapter Chapter

		page Page
		f    *os.File
		es   *elastic.Client
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

	if *saveJSON == true {
		f, err = os.Create("json/" + newid + ".json")
		if err != nil {
			return pages, err
		}

		defer f.Close()
	}

	if *indexDB == true {
		es, err = elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL("http://localhost:32769"),
		)
		if err != nil {
			log.Println(err)
			return pages, err
		}
	}

	count := 0

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
			PageIDValue, err = pageid.Value()
			if err != nil {
				log.Println(err)
				return pages, err
			}
		} else {
			PageIDValue = "null"
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
			PageID:     PageIDValue.(string),
			PageBody:   PageBodyValue.(string),
			Volume:     VolumeValue.(string),
			PageNumber: PageNumberValue.(string),
			Chapter:    chapter,
			Book:       book,
		}

		pages := append(pages, page)

		chapter, err = getChapter(db, id, page.PageNumber)
		if err != nil {
			log.Println(err)
			return pages, err
		}

		if *saveJSON == true {

			jsonByte, err := json.Marshal(page)
			if err != nil {
				log.Println(err)
				return pages, err
			}

			_, err = f.Write([]byte(fmt.Sprintf("%s\n", string(jsonByte))))
			if err != nil {
				log.Println(err)
				return pages, err
			}
		}

		if *indexDB == true {

			// index each page
			_, err := es.Index().Pretty(true).
				OpType("create").
				Index("maktabah").
				Type("pages").
				Id(newid + "-" + strconv.Itoa(count)).
				BodyJson(page).
				Do()

			if err != nil {
				log.Println(err)
				return pages, err
			}
		}

		count++
	}

	return pages, nil
}
