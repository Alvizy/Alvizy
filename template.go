package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var tmpl *template.Template

var funcMap = template.FuncMap{
	"splitComma": func(s string) []string {
		var result []string
		for _, p := range strings.Split(s, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
		return result
	},
	"formatDate": func(t time.Time) string {
		return t.Format("Jan 2006")
	},
	"adminPath": func() string {
		return "/" + cfg.AdminPath
	},
	"initial": func(s string) string {
		if len(s) == 0 {
			return ""
		}
		return string([]rune(s)[0])
	},
	"add": func(a, b int) int {
		return a + b
	},
}

func loadTemplates() {
	var err error
	tmpl, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("gagal load templates: %v", err)
	}
	// load admin templates ke set yang sama
	adminTmpl, err := tmpl.ParseGlob("templates/admin/*.html")
	if err == nil {
		tmpl = adminTmpl
	}
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
