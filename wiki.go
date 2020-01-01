package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
)

type Page struct {
    Title string
    Body []byte
}

func (p *Page) save() error {
    fname := p.Title + ".txt"
    return ioutil.WriteFile(fname, p.Body, 0600)
}

func load(title string) *Page {
    fname := title + ".txt"
    body, _ := ioutil.ReadFile(fname)
    return &Page{Title: title, Body: body}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p := load(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, string(p.Body))
}

func main_1() {
    p1 := Page{Title: "TestPage", Body: []byte("This is an example page")}
    p1.save()

    p2 := load("TestPage")
    fmt.Println(string(p2.Body))
}

func main() {
    http.HandleFunc("/view/", viewHandler)
    log.Fatal(http.ListenAndServe(":9090", nil))
}
