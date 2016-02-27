package main

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

func (s *Status) Poll() (bool, error) {

	// Put code here for checking file if indeed there are new urls

	return true, nil
}
