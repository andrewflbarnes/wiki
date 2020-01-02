package main

import (
    "io/ioutil"
    "net/http"
    "log"
    "html/template"
    "regexp"
)

var templates = template.Must(template.ParseGlob("tmpl/*"))
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }

        fn(w, r, m[2])
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := load(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}

    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":9090", nil))
}
