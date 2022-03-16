package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Data struct {
	News []string
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fileData, _ := ioutil.ReadFile("./index.tmpl")

	tmpl, err := template.New("index").Parse(string(fileData))
	if err != nil {
		fmt.Println(err)
	}

	data := Data{
		News: []string{
			"Current JIRA In Progress: POPS-2538",
			"No Blockers",
			"Next OOO: April 1-3",
		},
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", serveFiles)
	log.Fatal(http.ListenAndServe(":80", nil))
}
