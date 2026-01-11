package main

import (
	"html/template"
	"net/http"

	internaldb "forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/repo"
)

func main() {
	db := internaldb.InitDB()

	if err := repo.SeedCategories(db); err != nil {
		panic(err)
	}

	tpl := template.Must(template.ParseFiles("templates/layout.html"))
	store := repo.NewStore(db)
	app := &handlers.App{
		DB:       db,
		Tpl:      tpl,
		Posts:    store,
		Users:    store,
		Comments: store,
	}

	http.HandleFunc("/", app.HomeHandler)
	http.HandleFunc("/post", app.PostPageHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/logout", app.LogoutHandler)
	http.HandleFunc("/create-post", app.CreatePostHandler)
	http.HandleFunc("/addcomment", app.CommentHandler)
	http.HandleFunc("/react-post", app.ReactPosts)
	http.HandleFunc("/react-comment", app.ReactComment)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
