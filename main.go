package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
)

type Page struct {
	// Title    string
	// Template string
	Content template.HTML
	Path    string
	Title   string
}

func fatalOut(message string, params ...interface{}) {
	red := color.New(color.FgRed)
	_, _ = red.Printf("Fatal: "+message, params...)
	os.Exit(1)
}

func warnOut(message string, params ...interface{}) {
	red := color.New(color.FgHiYellow)
	_, _ = red.Printf("Warning: "+message, params...)
}

func loadPage(filename string) *Page {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		fatalOut("Error reading page file: %v\n", err)
	}
	path := strings.Replace(filename, ".md", ".html", 1)
	content := markdown.ToHTML(md, nil, nil)
	page := Page{template.HTML(string(content)), path, "My Title Here"}
	return &page
}

func loadPages(dir string) []*Page {

	var pages []*Page

	c, err := ioutil.ReadDir(dir)
	if err != nil {
		fatalOut("Error reading directory %v.\n", dir)
	}
	for _, v := range c {
		if v.IsDir() {
			pages = append(pages, loadPages(dir+"/"+v.Name())...)
		} else if v.Name() == "index.md" {
			pages = append(pages, loadPage(dir+"/"+v.Name()))
		} else {
			warnOut("Unknown file %v.\n", dir)
		}
	}
	return pages
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
			// todo build out init feature
			return
		} else {
			fatalOut("Unknown command %v\n", args[0])
		}
	}

	// if we get here the first param was a path
	inDir := args[0]
	// outDir := "www" // todo hardcoded for now

	// load each page of the site
	for _, page := range loadPages(inDir + "/pages") {
		tmpl, err := template.ParseFiles(inDir + "/templates/index.html")
		if err != nil {
			fatalOut("Error parsing template '%v' - %v.\n", "templates/index.html", err)
		}
		// tmpl.Execute(outDir+"/"+page.Path, page)
		tmpl.Execute(os.Stdout, page)
	}

	// load templates
	c, err := ioutil.ReadDir(inDir + "/templates")
	if err != nil {
		fatalOut("A '%vtemplates/' directory must exist to contain site templates.\n", inDir)
	}
	if len(c) == 0 {
		fatalOut("Found no templates in '%vtemplates' directory.\n", inDir)
	}

	files := make([]string, len(c))
	for i, v := range c {
		if strings.LastIndex(v.Name(), ".html") == len(v.Name())-5 {
			files[i] = inDir + "templates/" + v.Name()
		}
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		fatalOut("Error parsing template %v\n", err)
	}

	fmt.Println(templates)
}
