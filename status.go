package main

import (
	"bufio"
	"fmt"
	"os"
)

// poll json file checking for new urls and if such new urls exist, extract the url
// and pass it on to the services responsible for:
//
// 1. Download from shamela
// 2. Extracting data from .bok filetype
// 3. Indexing into Elasticsearch

type Status struct {

	// Responsible for returning a boolean if indeed there are new urls
	Exists bool
}

func (s *Status) Poll() error {
	// TODO: poll file and see if there any new lines, and return status

	// open file for reading
	f, err := os.Open("urls.json")
	defer f.Close()
	//f, err := ioutil.ReadFile("urls.json")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		fmt.Println(scanner.Text())
	}
	return nil
}
