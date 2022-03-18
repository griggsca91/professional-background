package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Data struct {
	News []string
}

func readData() []string {
	data, err := ioutil.ReadFile("feeds.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	filteredLines := make([]string, 0)
	for _, line := range lines {
		cleanedLine := strings.TrimSpace(line)
		if len(cleanedLine) > 0 {
			filteredLines = append(filteredLines, cleanedLine)
		}
	}

	return filteredLines
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fileData, _ := ioutil.ReadFile("./index.tmpl")

	tmpl, err := template.New("index").Parse(string(fileData))
	if err != nil {
		fmt.Println(err)
	}

	news := readData()

	data := Data{
		News: news,
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", serveFiles)
	log.Fatal(http.ListenAndServe(":80", nil))
}
