package main

import (
    "io/ioutil"
    "net/http"
    "log"
    "html/template"
)

type Page struct {
    Title string
    Body []byte
}

func (p *Page) save() error {
    fname := p.Title + ".txt"
    return ioutil.WriteFile(fname, p.Body, 0600)
}

func load(title string) (*Page, error) {
    fname := title + ".txt"
    body, e := ioutil.ReadFile(fname)
    if e == nil {
        return &Page{Title: title, Body: body}, e
    } else {
        return nil, e
    }
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := load(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    log.Fatal(http.ListenAndServe(":9090", nil))
}
