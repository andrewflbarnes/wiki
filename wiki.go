package main

import (
    "io/ioutil"
    "net/http"
    "log"
    "html/template"
    "regexp"
    "errors"
)

var templates = template.Must(template.ParseFiles("view.html", "edit.html"))
var validPath = regexp.MustCompile("^/(edit|view|save)/([0-9a-zA-Z]+)$")

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

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }

    return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := load(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}

    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":9090", nil))
}
