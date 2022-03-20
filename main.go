package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/griggsca91/professionalbackground/api"
	"github.com/spf13/viper"
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

	client := api.NewJiraClient(
		viper.Get("siteURL").(string),
		viper.Get("userEmail").(string),
		viper.Get("accessToken").(string),
	)

	fmt.Println(r.URL.Path)
	fileData, _ := ioutil.ReadFile("./index.tmpl")

	tmpl, err := template.New("index").Parse(string(fileData))
	if err != nil {
		fmt.Println(err)
	}

	news := readData()

	ticket, _ := client.GetLatestBlockedTicket()
	if ticket != "" {
		news = append(news, fmt.Sprintf("Blocked Ticket: %s", ticket))
	}

	ticket, _ = client.GetLatestInProgressTicket()
	if ticket != "" {
		news = append(news, fmt.Sprintf("Current ticket in progress: %s", ticket))
	}

	data := Data{
		News: news,
	}

	tmpl.Execute(w, data)
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	http.HandleFunc("/", serveFiles)
	log.Fatal(http.ListenAndServe(":80", nil))
}
