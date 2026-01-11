package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/internal/middleware"
)

func (a *App) ReactPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}

	user, err := middleware.CurrentUser(a.DB, r)
	if err != nil {
		a.renderError(w, http.StatusUnauthorized, "Вы должны авторизоваться, чтобы ставить лайки", nil)
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
	valueStr := r.FormValue("value")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		a.renderError(w, http.StatusBadRequest, "Некорректный post_id", user)
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

	value, err := strconv.Atoi(valueStr)
	if err != nil || (value != 1 && value != -1) {
		a.renderError(w, http.StatusBadRequest, "Некорректное значение реакции", user)
		return
	}

	var existing int
	row := a.DB.QueryRow(`SELECT value FROM post_reactions WHERE user_id = ? AND post_id = ?`, user.ID, postID)

	err = row.Scan(&existing)

	if err == sql.ErrNoRows {
		_, err = a.DB.Exec(`INSERT INTO post_reactions (user_id, post_id, value) VALUES (?, ?, ?)`,
			user.ID, postID, value)

		if err != nil {
			a.logError(err, "insert post reaction")
			a.renderError(w, http.StatusInternalServerError, "Ошибка сохранения реакции", user)
			return
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	if err != nil {
		a.logError(err, "select post reaction")
		a.renderError(w, http.StatusInternalServerError, "Ошибка чтения реакции", user)
		return
	}

	if existing == value {
		_, err = a.DB.Exec(
			`DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?`,
			user.ID, postID,
		)
		if err != nil {
			a.logError(err, "delete post reaction")
			a.renderError(w, http.StatusInternalServerError, "Ошибка удаления реакции", user)
			return
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	_, err = a.DB.Exec(
		`UPDATE post_reactions SET value = ? WHERE user_id = ? AND post_id = ?`,
		value, user.ID, postID,
	)
	if err != nil {
		a.logError(err, "update post reaction")
		a.renderError(w, http.StatusInternalServerError, "Ошибка обновления реакции", user)
		return
	}

	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (a *App) ReactComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}

	user, err := middleware.CurrentUser(a.DB, r)
	if err != nil {
		a.renderError(w, http.StatusUnauthorized, "Нужно войти", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.renderError(w, http.StatusBadRequest, "Некорректная форма", user)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	valueStr := r.FormValue("value")
	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		a.renderError(w, http.StatusBadRequest, "Некорректный comment_id", user)
		return
	}
	exists, err := a.Comments.CommentExists(commentID)
	if err != nil {
		a.logError(err, "comment exists")
		a.renderError(w, http.StatusInternalServerError, "Ошибка проверки комментария", user)
		return
	}
	if !exists {
		a.renderError(w, http.StatusNotFound, "Комментарий не найден", user)
		return
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || (value != 1 && value != -1) {
		a.renderError(w, http.StatusBadRequest, "Некорректное значение реакции", user)
		return
	}

	var existing int
	row := a.DB.QueryRow(
		`SELECT value FROM comment_reactions WHERE user_id = ? AND comment_id = ?`,
		user.ID, commentID,
	)
	err = row.Scan(&existing)

	if err == sql.ErrNoRows {
		_, err = a.DB.Exec(
			`INSERT INTO comment_reactions (user_id, comment_id, value) VALUES (?, ?, ?)`,
			user.ID, commentID, value,
		)
		if err != nil {
			a.logError(err, "insert comment reaction")
			a.renderError(w, http.StatusInternalServerError, "Ошибка сохранения реакции", user)
			return
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	if err != nil {
		a.logError(err, "select comment reaction")
		a.renderError(w, http.StatusInternalServerError, "Ошибка чтения реакции", user)
		return
	}

	if existing == value {
		_, err = a.DB.Exec(
			`DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?`,
			user.ID, commentID,
		)
		if err != nil {
			a.logError(err, "delete comment reaction")
			a.renderError(w, http.StatusInternalServerError, "Ошибка удаления реакции", user)
			return
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	_, err = a.DB.Exec(
		`UPDATE comment_reactions SET value = ? WHERE user_id = ? AND comment_id = ?`,
		value, user.ID, commentID,
	)
	if err != nil {
		a.logError(err, "update comment reaction")
		a.renderError(w, http.StatusInternalServerError, "Ошибка обновления реакции", user)
		return
	}

	http.Redirect(w, r, next, http.StatusSeeOther)
}
