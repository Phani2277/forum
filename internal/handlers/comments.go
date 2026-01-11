package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/middleware"
)

func (a *App) CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}

	user, err := middleware.CurrentUser(a.DB, r)
	if err != nil {
		a.renderError(w, http.StatusUnauthorized, "Вы должны авторизоваться, чтобы комментировать", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.renderError(w, http.StatusBadRequest, "Некорректная форма", user)
		return
	}
	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}
	postIDStr := r.FormValue("post_id")
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		a.renderError(w, http.StatusBadRequest, "Комментарий не может быть пустым", user)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		a.renderError(w, http.StatusBadRequest, "Неверный post_id", user)
		return
	}
	exists, err := a.Posts.PostExists(postID)
	if err != nil {
		a.logError(err, "post exists")
		a.renderError(w, http.StatusInternalServerError, "Ошибка проверки поста", user)
		return
	}
	if !exists {
		a.renderError(w, http.StatusNotFound, "Пост не найден", user)
		return
	}

	if err := a.Comments.CreateComment(postID, user.ID, content); err != nil {
		a.logError(err, "create comment")
		a.renderError(w, http.StatusInternalServerError, "Ошибка при создании комментария", user)
		return
	}

	http.Redirect(w, r, next, http.StatusSeeOther)
}
