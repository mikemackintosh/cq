package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/gosimple/slug"
)

var (
	// flagFormat  string
	// flagFile    string
	// flagLineNum bool
	flagShowHeaders bool
)

func init() {
	// flag.StringVar(&flagFile, "f", "", "File to parse")
	// flag.BoolVar(&flagLineNum, "n", false, "Shows line numbers")
	flag.BoolVar(&flagShowHeaders, "headers", false, "Shows line headers")
	// flag.StringVar(&flagFormat, "o", "", "Output Format")
}

// Collection is a collection of rows.
type Collection struct {
	Header *Header
	Rows   map[int]Row
}

// NewCollection creates a new default collection.
func NewCollection() *Collection {
	return &Collection{
		Header: &Header{},
		Rows:   make(map[int]Row),
	}
}

// Formatter takes a string and parses the template
type Formater struct {
	Format string
	Vars   interface{}
}

// NewFormater xreates a new formatter
func NewFormater(format string) (*Formater, error) {
	re := regexp.MustCompile(`\\\((\..*?)\)`)
	format = re.ReplaceAllString(format, "{{ $1 }}")

	if string(format[len(format)-1]) != "\n" {
		format = format + "\n"
	}
	fmtr := &Formater{Format: format}
	return fmtr, nil
}

// Parse the formatter
func (fmtr *Formater) Parse(linenum int, fields map[string]interface{}) error {
	fm := template.FuncMap{
		"join": strings.Join,
		"replace": func(x, y, z string) string {
			return strings.Replace(z, x, y, -1)
		},
	}
	tmpl, err := template.New("main").Funcs(fm).Parse(fmtr.Format)
	if err != nil {
		return err
	}

	fields["LINENUM"] = linenum
	if err := tmpl.Execute(os.Stdout, fields); err != nil {
		return err
	}

	return nil
}

// Headers is a map of column name to column number.
type Header struct {
	Name     map[int]string
	SlugName map[int]string
	ID       map[string]int
}

// Get will get a column name from column number.
func (h Header) Get(i int) (string, error) {
	if v, ok := h.SlugName[i]; ok {
		return v, nil
	}

	// Column not found, throw an error.
	return "", fmt.Errorf("could not find column")
}

// Row is a list of columns to values as an Interface.
type Row map[string]interface{}

// perror takes an error and if not nil, prints it and exits.
func perror(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// main parses the flags (if provided) and performs the script.
func main() {
	flag.Parse()

	// Check for arg counts
	if len(os.Args[:len(os.Args)]) < 2 {
		perror(errors.New("Please provide an output format"))
	}

	// Supply format string to formater
	fm, _ := NewFormater(os.Args[len(os.Args)-1])

	// Get stdin
	input := os.Stdin
	rows, err := csv.NewReader(input).ReadAll()
	if err != nil {
		perror(fmt.Errorf("Error parsesing csv, %s", err))
	}

	// This is a subjective approach, but we want at least 1 row and 1 column.
	// Anything less than that doesn't make much sense.
	if err := validateInput(rows); err != nil {
		perror(err)
	}

	// Create a new collection
	var c = NewCollection()

	// Generate headers
	var titleRow = 0
	if err := c.parseHeader(rows[titleRow]); err != nil {
		perror(err)
	}

	// Only show headers if we set it
	if flagShowHeaders {
		fmt.Printf("This spreadsheet has the following columns:\n\n")
		for id, name := range c.Header.SlugName {
			fmt.Printf("- %d, %s - Use: '\\(.%s)'\n", id, c.Header.Name[id], name)
		}
		os.Exit(0)
	}

	// Parse the rows
	if err := c.parseRows(rows[titleRow+1:]); err != nil {
		perror(err)
	}

	// Parse and Print
	for i, row := range c.Rows {
		// Parse lines
		if err = fm.Parse(i, row); err != nil {
			perror(err)
		}
	}
}

// parseHeader will parse the first row for column titles.
func (c *Collection) parseHeader(titleRow []string) error {
	var header = &Header{
		Name:     map[int]string{},
		SlugName: map[int]string{},
		ID:       map[string]int{},
	}

	// loop through the title row and create the mappings for later.
	for i, v := range titleRow {
		header.Name[i] = v
		header.SlugName[i] = strings.Replace(slug.Make(v), "-", "_", -1)
		header.ID[v] = i
	}

	c.Header = header
	return nil
}

// parseRows parses the rows and populates the array.
func (c *Collection) parseRows(rows [][]string) error {
	for rowNum, columns := range rows {
		var row = map[string]interface{}{}
		for columnnum, value := range columns {
			columnName, err := c.Header.Get(columnnum)
			if err != nil {
				return err
			}
			row[columnName] = value
		}

		c.Rows[rowNum+1] = row
	}

	return nil
}

// validateInput just checks to make sure it's not empty or too small of a csv.
func validateInput(rows [][]string) error {
	if len(rows) > 0 {
		if len(rows[0]) > 0 {
			return nil
		}
	}

	return errors.New("your csv should have columns and rows")
}
