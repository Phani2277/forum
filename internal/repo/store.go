package repo

import (
	"database/sql"

	"forum/internal/models"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateUser(email string, username string, password string) error {
	return CreateUser(s.DB, email, username, password)
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	return GetUserByEmail(s.DB, email)
}

func (s *Store) GetUserByID(id int) (*models.User, error) {
	return GetUserByID(s.DB, id)
}

func (s *Store) CreatePost(userID int, title string, content string, categoryIDs []int) (int, error) {
	return CreatePost(s.DB, userID, title, content, categoryIDs)
}

func (s *Store) GetPostCards(filter PostCardsFilter) ([]models.PostCard, error) {
	return GetPostCards(s.DB, filter)
}

func (s *Store) GetPostCardWithComments(postID int) (*models.PostCardWithComments, error) {
	return GetPostCardWithComments(s.DB, postID)
}

func (s *Store) PostExists(postID int) (bool, error) {
	return PostExists(s.DB, postID)
}

func (s *Store) CreateComment(postID int, userID int, content string) error {
	return CreateComment(s.DB, postID, userID, content)
}

func (s *Store) GetCommentsByPostID(postID int) ([]models.CommentView, error) {
	return GetCommentsByPostID(s.DB, postID)
}

func (s *Store) CommentExists(commentID int) (bool, error) {
	return CommentExists(s.DB, commentID)
}
