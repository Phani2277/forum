package handlers

import (
	"net/http"
	"strings"
	"time"

	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user, _ := middleware.CurrentUser(a.DB, r)
		data := models.BasePageData{CurrentUser: user}
		a.render(w, "register.html", data)
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data := models.BasePageData{Error: "Некорректная форма"}
			a.renderWithStatus(w, http.StatusBadRequest, "register.html", data)
			return
		}
		email := strings.TrimSpace(r.FormValue("email"))
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		if email == "" || username == "" || password == "" {
			data := models.BasePageData{Error: "Заполните email, username и password"}
			a.renderWithStatus(w, http.StatusBadRequest, "register.html", data)
			return
		}

		err := a.Users.CreateUser(email, username, password)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				data := models.BasePageData{Error: "Пользователь с таким email уже существует"}
				a.renderWithStatus(w, http.StatusBadRequest, "register.html", data)
				return
			}
			a.logError(err, "create user")
			data := models.BasePageData{Error: "Ошибка регистрации"}
			a.renderWithStatus(w, http.StatusInternalServerError, "register.html", data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user, _ := middleware.CurrentUser(a.DB, r)
		data := models.BasePageData{CurrentUser: user}
		a.render(w, "login.html", data)
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data := models.BasePageData{Error: "Некорректная форма"}
			a.renderWithStatus(w, http.StatusBadRequest, "login.html", data)
			return
		}
		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))
		if email == "" || password == "" {
			data := models.BasePageData{Error: "Введите email и password"}
			a.renderWithStatus(w, http.StatusBadRequest, "login.html", data)
			return
		}

		user, err := a.Users.GetUserByEmail(email)
		if err != nil {
			a.logError(err, "get user by email")
			data := models.BasePageData{Error: "Пользователь не найден"}
			a.renderWithStatus(w, http.StatusNotFound, "login.html", data)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			data := models.BasePageData{Error: "Пароль неверный"}
			a.renderWithStatus(w, http.StatusUnauthorized, "login.html", data)
			return
		}

		sessionID, err := repo.CreateSessions(a.DB, user.ID)
		if err != nil {
			a.logError(err, "create session")
			a.renderError(w, http.StatusInternalServerError, "Ошибка сессии", nil)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   sessionID,
			Expires: time.Now().Add(20 * time.Minute),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
}

func (a *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}
	c, err := r.Cookie("session")
	if err == nil {
		_ = repo.DeleteSession(a.DB, c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		Path:    "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
