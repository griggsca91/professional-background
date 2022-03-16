package main

import "fmt"
import "net/http"
import "log"


func serveFiles(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Path)
    p := "." + r.URL.Path
    if p == "./" {
        p = "./index.html"
    }
    http.ServeFile(w, r, p)
}


func main() {
      http.HandleFunc("/", serveFiles)
    log.Fatal(http.ListenAndServe(":80", nil))


}
