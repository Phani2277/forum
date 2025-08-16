package server

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/models"
	"forum/internal/session"
)

// App holds application dependencies.
type App struct {
	DB        *sql.DB
	Templates *template.Template
}

// New creates a new App instance.
func New(db *sql.DB) *App {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &App{DB: db, Templates: tmpl}
}

// Helper to get current user from request.
func (a *App) currentUser(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	userID, err := session.GetUserID(a.DB, cookie.Value)
	if err != nil {
		return nil
	}
	u, err := models.GetByID(a.DB, userID)
	if err != nil {
		return nil
	}
	return u
}

// Home handler displays posts.
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	filter := r.URL.Query().Get("filter")
	var userID, likedBy int
	cu := a.currentUser(r)
	if filter == "my" && cu != nil {
		userID = cu.ID
	}
	if filter == "liked" && cu != nil {
		likedBy = cu.ID
	}
	posts, err := models.GetPosts(a.DB, category, userID, likedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	categories, _ := models.GetCategories(a.DB)
	data := map[string]interface{}{
		"Posts":       posts,
		"Categories":  categories,
		"CurrentUser": cu,
		"Filter":      filter,
		"Category":    category,
	}
	a.Templates.ExecuteTemplate(w, "index.html", data)
}

// Register displays and handles user registration.
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.Templates.ExecuteTemplate(w, "register.html", nil)
	case http.MethodPost:
		email := strings.TrimSpace(r.FormValue("email"))
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		if email == "" || username == "" || password == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		if err := models.CreateUser(a.DB, email, username, password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Login displays and handles login.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.Templates.ExecuteTemplate(w, "login.html", nil)
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := models.Authenticate(a.DB, email, password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := session.Create(a.DB, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: token, Path: "/", Expires: time.Now().Add(session.Duration)})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Logout terminates user session.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err == nil {
		session.Delete(a.DB, cookie.Value)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreatePost displays form and handles new post creation.
func (a *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	cu := a.currentUser(r)
	if cu == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		categories, _ := models.GetCategories(a.DB)
		data := map[string]interface{}{"Categories": categories}
		a.Templates.ExecuteTemplate(w, "create_post.html", data)
	case http.MethodPost:
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		cats := strings.Split(r.FormValue("categories"), ",")
		for i := range cats {
			cats[i] = strings.TrimSpace(cats[i])
		}
		if title == "" || content == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		id, err := models.CreatePost(a.DB, cu.ID, title, content, cats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(id), http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ViewPost displays a post and handles comments.
func (a *App) ViewPost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		post, err := models.GetPost(a.DB, id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		comments, _ := models.GetComments(a.DB, id)
		data := map[string]interface{}{
			"Post":        post,
			"Comments":    comments,
			"CurrentUser": a.currentUser(r),
		}
		a.Templates.ExecuteTemplate(w, "post.html", data)
	case http.MethodPost:
		cu := a.currentUser(r)
		if cu == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Error(w, "empty comment", http.StatusBadRequest)
			return
		}
		if err := models.CreateComment(a.DB, id, cu.ID, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// LikePost handles liking/disliking posts.
func (a *App) LikePost(w http.ResponseWriter, r *http.Request) {
	cu := a.currentUser(r)
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.Split(idStr, "/")
	if len(parts) < 2 || parts[1] != "like" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	value, _ := strconv.Atoi(r.FormValue("value"))
	if value != 1 && value != -1 && value != 0 {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}
	if err := models.SetPostLike(a.DB, cu.ID, id, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(id), http.StatusSeeOther)
}

// LikeComment handles liking/disliking comments.
func (a *App) LikeComment(w http.ResponseWriter, r *http.Request) {
	cu := a.currentUser(r)
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/comment/")
	parts := strings.Split(idStr, "/")
	if len(parts) < 2 || parts[1] != "like" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	value, _ := strconv.Atoi(r.FormValue("value"))
	if value != 1 && value != -1 && value != 0 {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}
	if err := models.SetCommentLike(a.DB, cu.ID, id, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var postID int
	a.DB.QueryRow("SELECT post_id FROM comments WHERE id=?", id).Scan(&postID)
	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
