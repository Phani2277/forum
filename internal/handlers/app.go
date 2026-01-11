package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/repo"
)

type App struct {
	DB       *sql.DB
	Tpl      *template.Template
	Posts    PostRepo
	Users    UserRepo
	Comments CommentRepo
}

func (a *App) render(w http.ResponseWriter, page string, data any) {
	a.renderWithStatus(w, http.StatusOK, page, data)
}

func (a *App) renderWithStatus(w http.ResponseWriter, status int, page string, data any) {
	tmpl, err := a.Tpl.Clone()
	if err != nil {
		log.Printf("template clone error: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	if _, err := tmpl.ParseFiles("templates/" + page); err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) renderError(w http.ResponseWriter, status int, message string, user *models.User) {
	data := models.ErrorPageData{
		CurrentUser: user,
		Status:      status,
		Message:     message,
	}
	a.renderWithStatus(w, status, "error.html", data)
}

func (a *App) logError(err error, message string) {
	if err == nil {
		return
	}
	log.Printf("%s: %v", message, err)
}

type PostRepo interface {
	CreatePost(userID int, title string, content string, categoryIDs []int) (int, error)
	GetPostCards(filter repo.PostCardsFilter) ([]models.PostCard, error)
	GetPostCardWithComments(postID int) (*models.PostCardWithComments, error)
	PostExists(postID int) (bool, error)
}

type UserRepo interface {
	CreateUser(email string, username string, password string) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type CommentRepo interface {
	CreateComment(postID int, userID int, content string) error
	GetCommentsByPostID(postID int) ([]models.CommentView, error)
	CommentExists(commentID int) (bool, error)
}
