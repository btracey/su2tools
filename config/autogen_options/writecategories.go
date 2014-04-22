package main

import (
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func writeCategoriesAndOptionOrder(categories []*ConfigCategory, options []*ConfigOption) {

	categoryFilename := filepath.Join(configPath, "category_and_order.go")

	categoryFile, err := os.Create(categoryFilename)
	if err != nil {
		log.Fatalf("error creating %s: %v", categoryFilename, err)
	}

	categoryFile.WriteString(packageHeader)
	categoryFile.WriteString(autogenmessage)

	// First, print all of the categories, names, and descriptions
	categoryFile.WriteString(
		`type category struct{
	Id int
	Name string
	Description string
}

var categoryList = []category{
`)

	maxCat := len(categories)

	for _, cat := range categories {
		categoryFile.WriteString("{\n" +
			"Id: " + strconv.Itoa(cat.Id) + ",\n" +
			"Name: \"" + cat.Name + "\",\n" +
			"Description: \"" + cat.Description + "\",\n" +
			"},\n")
	}

	categoryFile.WriteString("}\n")
	// Now print the options in order

	categoryFile.WriteString(
		`
		// optionList contains a list of the options to be printed. Outer slice is
		// for the category, inner list is the order of the options
		var optionList = [][]Option{
		`)
	var cat = -1 // Category of the previous option

	for _, opt := range options {
		if opt.Category != cat {
			if cat != -1 {
				// Finish off the previos
				categoryFile.WriteString("},\n")
			}
			if cat != maxCat {
				// Start a new one
				categoryFile.WriteString("{\n")
			}
		}
		// Add a new option
		categoryFile.WriteString(opt.Value + ",\n")
		cat = opt.Category
	}

	categoryFile.WriteString("},\n}\n")

	categoryFile.Close()

	b, err := ioutil.ReadFile(categoryFilename)
	if err != nil {
		log.Fatalf("error creating %s: %v", categoryFilename, err)
	}
	b, err = format.Source(b)
	if err != nil {
		log.Fatalf("error formatting %s: %v", categoryFilename, err)
	}
	err = ioutil.WriteFile(categoryFilename, b, 0700)
	if err != nil {
		log.Fatalf("error writing %s: %v", categoryFilename, err)
	}

}
