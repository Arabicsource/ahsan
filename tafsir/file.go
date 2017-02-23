package main

import (
	"io/ioutil"
	"os"
)

// readFromFile ...
func readFromFile(filepath string) ([]string, error) {
	// TODO: require to parse file flag and store the urls in a slice of type URL

	// initialize variable urls which is a slice of URL
	var urls []string

	// open file and check if there is an error
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	// read content of entire file and check if there is an error
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// iterate over each line of the data
	for _, url := range data {
		urls = append(urls, string(url))
	}

	return urls, nil
}
