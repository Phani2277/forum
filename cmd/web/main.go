package main

import (
	"log"
	"net/http"
	"strings"

	"forum/internal/database"
	"forum/internal/server"
)

func main() {
	db, err := database.Open("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.Init(db); err != nil {
		log.Fatal(err)
	}

	app := server.New(db)
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/register", app.Register)
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/logout", app.Logout)
	mux.HandleFunc("/post/create", app.CreatePost)
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/like") {
			app.LikePost(w, r)
			return
		}
		app.ViewPost(w, r)
	})
	mux.HandleFunc("/comment/", app.LikeComment)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
