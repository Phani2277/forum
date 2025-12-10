package main

import (
	"database/sql"
	"strconv"
	"time"

	"net/http"
)

var db *sql.DB

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := CurrentUser(r)

	posts, err := GetAllPostsd(db)

	if err != nil {
		w.Write([]byte("ошибка получения постов"))
		return

	}

	w.Write([]byte("<html><body>"))
	if user == nil {
		w.Write([]byte("<p>Привет, гость!</p>"))
	} else {
		w.Write([]byte("<p>Привет, " + user.Username + "!</p>"))
	}

	w.Write([]byte("<h1>Posts</h1>"))

	for _, p := range posts {
		w.Write([]byte("<hr>"))
		w.Write([]byte("<h2>" + p.Title + "</h2>"))
		w.Write([]byte("<p>" + p.Content + "</p>"))
		w.Write([]byte("<p>Category: " + p.Category + "</p>"))

		w.Write([]byte(`
    	<form method="POST" action="/addcomment">
        <input type="hidden" name="post_id" value="...">
        <input type="text" name="content" placeholder="Ваш комментарий"><br>
        <button type="submit">Отправить</button>
    	</form>
		`))

	}

	w.Write([]byte("</body></html>"))
}

func CurrentUser(r *http.Request) (*User, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	sessionID := c.Value

	user, err := GetUserBySessionID(db, sessionID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte(`
    <html>
    <body>
        <h1>Register</h1>
        <form method="POST" action="/register">
            <input type="email" name="email" placeholder="Email"><br>
            <input type="text" name="username" placeholder="Username"><br>
            <input type="password" name="password" placeholder="Password"><br>
            <button type="submit">Register</button>
        </form>
    </body>
    </html>
	`))
		return
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := CreateUser(db, email, username, password)

		if err != nil {
			panic(err)
		} else {
			w.Write([]byte("пользователь добавлен"))
			return
		}

	} else {
		w.Write([]byte("404"))
		return
	}

}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("Method not allowed"))
		return
	}

	user, err := CurrentUser(r)
	if err != nil {
		w.Write([]byte("Вы должны авторизоваться, чтобы комментировать"))
		return
	}

	r.ParseForm()
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		w.Write([]byte("Неверный post_id"))
		return
	}

	err = CreateComment(db, postID, user.ID, content)
	if err != nil {
		w.Write([]byte("Ошибка при создании комментария"))
		return
	}

	// Пока можно просто ответить текстом
	w.Write([]byte("Комментарий добавлен"))
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user, err := CurrentUser(r)

	if err != nil {
		w.Write([]byte("need authorization for create post"))
		return
	}

	if r.Method == http.MethodGet {
		w.Write([]byte(`
            <html><body>
            <h1>Create Post</h1>
            <form method="POST" action="/create-post">
                <input type="text" name="title" placeholder="Title"><br>
                <textarea name="content" placeholder="Content"></textarea><br>
                <input type="text" name="category" placeholder="Category"><br>
                <button type="submit">Создать пост</button>
            </form>
            </body></html>
        `))
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		err := CreatePost(db, user.ID, title, content, category)
		if err != nil {
			w.Write([]byte("error create post"))
			return
		}

		w.Write([]byte("post create"))
		return
	}

	w.Write([]byte("404"))

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte(`
        <html>
        <body>
            <h1>Login</h1>
            <form method="POST" action="/login">
                <input type="email" name="email" placeholder="Email"><br>
                <input type="password" name="password" placeholder="Password"><br>
                <button type="submit">Login</button>
            </form>
        </body>
        </html>
    `))
		return
	} else if r.Method == http.MethodPost {

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := GetUserByEmail(db, email)

		if err != nil {
			w.Write([]byte("Пользователь не найден"))
			return
		}

		if user.Password != password {
			w.Write([]byte("Пароль неверный"))
			return
		}

		sessionID, err := CreateSessions(db, user.ID)

		if err != nil {
			w.Write([]byte("ошибка сессии "))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   sessionID,
			Expires: time.Now().Add(20 * time.Minute),
			Path:    "/",
		})

		w.Write([]byte("Вход выполнен"))
		return

	} else {
		w.Write([]byte("404"))
		return
	}

}

func main() {

	db = InitDB()
	_ = db // чтобы не ругался линтер

	// err := createUser(db, "Sewq@mail.ru", "Roma", "1234")

	// if err != nil {
	// 	fmt.Println("ошибка", err)
	// } else {
	// 	fmt.Println("пользователь есть")
	// }

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/addcomment", CommentHandler)

	http.ListenAndServe(":8080", nil)

}
