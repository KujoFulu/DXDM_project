package main

import (
	"fmt"
	"net/http"
	"html/template"
	"path"
)

const SiteAssets = "site"
const Templates = SiteAssets + "/templates"

func indexFunc(w http.ResponseWriter, r *http.Request){
	fmt.Println("method:", r.Method, "on", "index") // print request method
	if r.Method == "GET" {
		t, err := template.ParseFiles(path.Join(Templates, "index.html"))
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil) // we could pass a struct in to apply formatting if we wanted
	} /*else if r.Method == "POST" {
		servePostRequest(w, r)
	}*/
}


func main() {
	http.HandleFunc("/", indexFunc)
	fs := http.FileServer(http.Dir("./site"))
	http.NewServeMux().Handle("/site/", http.StripPrefix("/site/", fs))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}