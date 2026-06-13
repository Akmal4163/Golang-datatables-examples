package main

import (
	"net/http"
	"database/sql"
	"log"
	"html/template"
	"example/datatables/handlers"
	_ "modernc.org/sqlite"

)

type Data struct {
	Title string
}

func main() {
	db, err := sql.Open("sqlite", "data_barang.db")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer db.Close()

	datatableHandler := &handlers.DatatableHandlers{
		DB: db,
	}

	tmpl := template.Must(template.ParseFiles("templates/datatables.html"))

	http.HandleFunc("/api/barang", datatableHandler.GetData)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := Data{
			Title: "datatables example",
		}
		tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}