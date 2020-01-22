package main

import (
	"fmt"
	"github.com/fatih/color"
	"html/template"
	"io/ioutil"
	"os"
	"encoding/json"
	"strings"
)

func fatalOut(message string, params ...interface{}) {
	red := color.New(color.FgRed)
	_, _ = red.Printf(message, params...)
	os.Exit(1)
}

type Page struct {

}

func loadPage(filename string) *Page {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fatalOut("Error reading page file\n", err)
	}
	var page Page
	err = json.Unmarshal(body, page)
	if err != nil {
		fatalOut("Error unmarshalling page file\n", err)
	}
	return &page
}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fatalOut("Insufficient parameters passed to command. Example: ~$> gophergen <input-dir>\n")
		return
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		if args[0] == "init" {
			fmt.Println("Initialising project")
		} else {
			fatalOut("Unknown command %v\n", args[0])
		}
	}

	// if we get here the first param was a path
	inDir := args[0]

	c, err := ioutil.ReadDir(inDir + "/templates")
	if err != nil {
		fatalOut("A '%vtemplates/' directory must exist to contain site templates.\n", inDir)
	}
	if len(c) == 0 {
		fatalOut("Found no templates in '%vtemplates' directory.\n", inDir)
	}

	files := make([]string, len(c))
	for i, v := range c {
		if strings.LastIndex(v.Name(), ".html") == len(v.Name()) - 5 {
			files[i] = inDir + "templates/" + v.Name()
		}
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		fatalOut("Error parsing template %v\n", err)
	}

	fmt.Println(templates)
}
