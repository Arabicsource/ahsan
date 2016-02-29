package main

import (
	"errors"
	"log"
	"regexp"
)

// Tag represents the HTML element that will be parsed and pulled from
// html document
type Tag struct {
	// Name of the HTML element like 'div' and 'a'.
	Name string

	// slice of the class names, and it may be empty if the element has no classes.
	Class []string

	// String consisting of the ID of a particular element, to be used by the regex.
	Id string

	// Text inside the element
	Text string

	// the Regex for that particular element to be parsed.
	Regex *regexp.Regexp
}

//New inits a new Tag (Html Element)
func (t *Tag) New(name string) *Tag {

	t.Name = name
	return t
}

//Compile returns a regular expression, and a bool
func (t *Tag) Compile(name string) (*regexp.Regexp, bool) {

	var ok bool
	var re *regexp.Regexp

	// Checking the given element against pre-defined cases and executing relevant regex
	// to pull out that element from the HTML
	// TODO:: refactor the regex bit to its own method
	switch t.Name {
	default:
		ok = false
		re = nil

		return re, ok

	// Do likewise for other HTML elements when you refactor the code
	case "a":
		// Looking for just about any anchor element
		log.Println("Instantiating for anchor element!")
		re = regexp.MustCompile(`\/index.php\/category\/\d+`)
		ok = true
		return re, ok
	}

}

// Match will match up a regular expression with a body of string (html document)
func (t *Tag) Match(re *regexp.Regexp, body string) (result []string, err error) {

	// TODO: Refactor the code so that the regex matches an entire element
	// and returns a list of elements.

	// Storing the list of books' urls in Crawler.Books
	result = re.FindAllString(string(body), -1)
	if result == nil {
		err := errors.New("Found no Match!")
		return nil, err
	}
	return result, nil
}
