package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/repo"
)

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.renderError(w, http.StatusMethodNotAllowed, "Метод не поддерживается", nil)
		return
	}
	user, _ := middleware.CurrentUser(a.DB, r)

	cats, err := repo.GetAllCategories(a.DB)
	if err != nil {
		a.logError(err, "get categories")
		a.renderError(w, http.StatusInternalServerError, "Ошибка категорий", user)
		return
	}

	liked := r.URL.Query().Get("liked")
	mine := r.URL.Query().Get("mine")
	categoryIDStr := r.URL.Query().Get("category_id")
	allActive := false
	mineActive := false
	likedActive := false
	selectedCategoryID := 0
	filter := repo.PostCardsFilter{
		CommentLimit: 3,
	}
	if liked == "1" {
		likedActive = true
		if user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		filter.LikedOnly = true
		filter.UserID = user.ID
	} else if mine == "1" {
		mineActive = true
		if user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		filter.MineOnly = true
		filter.UserID = user.ID
	} else if categoryIDStr != "" {
		categoryID, convErr := strconv.Atoi(categoryIDStr)
		if convErr != nil {
			a.renderError(w, http.StatusBadRequest, "Неверная категория", user)
			return
		}
		categoryFound := false
		for _, c := range cats {
			if c.ID == categoryID {
				categoryFound = true
				break
			}
		}
		if !categoryFound {
			a.renderError(w, http.StatusNotFound, "Категория не найдена", user)
			return
		}
		selectedCategoryID = categoryID
		filter.CategoryID = categoryID
	} else {
		allActive = true
	}

	cards, err := a.Posts.GetPostCards(filter)
	if err != nil {
		a.logError(err, "get post cards")
		a.renderError(w, http.StatusInternalServerError, "Ошибка получения постов", user)
		return
	}

	data := models.HomePageData{
		CurrentUser:        user,
		Categories:         cats,
		Posts:              cards,
		AllActive:          allActive,
		MineActive:         mineActive,
		LikedActive:        likedActive,
		SelectedCategoryID: selectedCategoryID,
	}

	a.render(w, "home.html", data)
}
