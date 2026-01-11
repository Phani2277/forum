package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/repo"
)

func (a *App) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.CurrentUser(a.DB, r)
	if err != nil {
		a.renderError(w, http.StatusUnauthorized, "Нужна авторизация для создания поста", nil)
		return
	}

	if r.Method == http.MethodGet {
		cats, err := repo.GetAllCategories(a.DB)
		if err != nil {
			a.logError(err, "get categories")
			a.renderError(w, http.StatusInternalServerError, "Ошибка категорий", user)
			return
		}
		data := models.CreatePostPageData{
			CurrentUser: user,
			Categories:  cats,
		}
		a.render(w, "create_post.html", data)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			cats, _ := repo.GetAllCategories(a.DB)
			data := models.CreatePostPageData{
				CurrentUser: user,
				Categories:  cats,
				Error:       "Некорректная форма",
			}
			a.renderWithStatus(w, http.StatusBadRequest, "create_post.html", data)
			return
		}
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		if title == "" || content == "" {
			cats, _ := repo.GetAllCategories(a.DB)
			data := models.CreatePostPageData{
				CurrentUser: user,
				Categories:  cats,
				Error:       "Заполните заголовок и текст",
			}
			a.renderWithStatus(w, http.StatusBadRequest, "create_post.html", data)
			return
		}
		catIDStrs := r.Form["category_id"]
		if len(catIDStrs) == 0 {
			cats, _ := repo.GetAllCategories(a.DB)
			data := models.CreatePostPageData{
				CurrentUser: user,
				Categories:  cats,
				Error:       "Выберите хотя бы одну категорию",
			}
			a.renderWithStatus(w, http.StatusBadRequest, "create_post.html", data)
			return
		}

		cats, err := repo.GetAllCategories(a.DB)
		if err != nil {
			a.logError(err, "get categories")
			a.renderError(w, http.StatusInternalServerError, "Ошибка категорий", user)
			return
		}
		validCats := make(map[int]bool, len(cats))
		for _, c := range cats {
			validCats[c.ID] = true
		}

		categoryIDs := make([]int, 0, len(catIDStrs))
		for _, catIDStr := range catIDStrs {
			catID, err := strconv.Atoi(catIDStr)
			if err != nil {
				data := models.CreatePostPageData{
					CurrentUser: user,
					Categories:  cats,
					Error:       "Неверная категория",
				}
				a.renderWithStatus(w, http.StatusBadRequest, "create_post.html", data)
				return
			}
			if !validCats[catID] {
				data := models.CreatePostPageData{
					CurrentUser: user,
					Categories:  cats,
					Error:       "Категория не найдена",
				}
				a.renderWithStatus(w, http.StatusNotFound, "create_post.html", data)
				return
			}
			categoryIDs = append(categoryIDs, catID)
		}

		if _, err := a.Posts.CreatePost(user.ID, title, content, categoryIDs); err != nil {
			a.logError(err, "create post")
			data := models.CreatePostPageData{
				CurrentUser: user,
				Categories:  cats,
				Error:       "Ошибка создания поста",
			}
			a.renderWithStatus(w, http.StatusInternalServerError, "create_post.html", data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", user)
}

func (a *App) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}
	user, _ := middleware.CurrentUser(a.DB, r)

	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		a.renderError(w, http.StatusBadRequest, "Некорректный id поста", user)
		return
	}

	post, err := a.Posts.GetPostCardWithComments(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.renderError(w, http.StatusNotFound, "Пост не найден", user)
			return
		}
		a.logError(err, "get post")
		a.renderError(w, http.StatusInternalServerError, "Ошибка загрузки поста", user)
		return
	}

	data := models.PostPageData{
		CurrentUser: user,
		Post:        *post,
	}

	a.render(w, "post.html", data)
}
