package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tapaap/config"
	"tapaap/utils"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ()
func main() {
	connStr := "host=localhost port=5433 user=myuser password=mypass dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT
	)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'active', -- или 'active', 'done'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
	)`)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlBytes, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			return
		}
		html := string(htmlBytes)
		cookie, err := r.Cookie("email")
		if err == nil {
			html = strings.Replace(html, `<div id="greeting"></div>`, `<div>Привет, `+cookie.Value+`</div>`, 1)
		}

		fmt.Fprint(w, html)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "about.html")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			showAboutWithError(w, "Все поля обязательны")
			return
		}
		var user struct {
			ID       int
			Password string
			Name     string
			Email    string
		}
		err := db.QueryRow(`SELECT id, password, name, email FROM users where email =$1`, email).
			Scan(&user.ID, &user.Password, &user.Name, &user.Email)
		if err != nil {
			showAboutWithError(w, "Неверная почта")
			return
		}
		if !utils.CheckPasswordHash(password, user.Password) {
			showAboutWithError(w, "Неверный пароль")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    user.Email,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   86400,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/registerURL", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "register.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			showAboutWithError(w, "Все поля обязательны")
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			showRegisterWithError(w, "Ошибка сервера")
			return
		}

		_, err = db.Exec(`INSERT INTO users (email, password) VALUES ($1, $2)`,
			email, hashedPassword)

		if err != nil {
			// Если email уже существует — будет ошибка UNIQUE
			showRegisterWithError(w, "Пользователь с такой почтой уже существует")
			return
		}

		var NewUser struct {
			Email string
		}

		err = db.QueryRow(`SELECT email FROM users where email = $1`, email).Scan(&NewUser.Email)
		if err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "email",
				Value:    NewUser.Email,
				HttpOnly: true,
				Path:     "/",
				MaxAge:   86400,
			})
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/auth/yandex/login", func(w http.ResponseWriter, r *http.Request) {
		authURL := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=login:email",
			config.YandexClientID,
			url.QueryEscape(config.YandexRedirectURI),
		)
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	http.HandleFunc("/auth/yandex/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Fprintln(w, "код яндекс пользователя не получен")
			return
		}
		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("client_id", config.YandexClientID)
		data.Set("client_secret", config.YandexClientSecret)

		resp, err := http.PostForm("https://oauth.yandex.ru/token", data)
		if err != nil {
			fmt.Fprint(w, "Ошибка при обмене кода на токен:", err)
			return
		}
		defer resp.Body.Close()

		var tokenResp struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			fmt.Fprintln(w, "Ошибка парсинга токена:", err)
			return
		}
		userInfoURL := fmt.Sprintf("https://login.yandex.ru/info?format=json&oauth_token=%s",
			tokenResp.AccessToken,
		)

		userResp, err := http.Get(userInfoURL)
		if err != nil {
			fmt.Fprintln(w, "ошибка получения данных пользователя:", err)
			return
		}
		defer userResp.Body.Close()

		var yandexUser struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			DefaultEmail string `json:"default_email"`
		}

		if err := json.NewDecoder(userResp.Body).Decode(&yandexUser); err != nil {
			fmt.Fprintln(w, "Ошибка парсинга пользователя:", err)
			return
		}

		err = db.QueryRow(`SELECT email FROM users where email = $1`, yandexUser.DefaultEmail).Scan(new(string))
		if err == sql.ErrNoRows {
			_, err = db.Exec(`INSERT INTO users (email, password, name) VALUES ($1, $2, $3)`,
				yandexUser.DefaultEmail, "", yandexUser.FirstName+" "+yandexUser.LastName)
			if err != nil {
				http.Error(w, "Ошибка регистрации через Яндекс", http.StatusInternalServerError)
				return
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    yandexUser.DefaultEmail,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   86400,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("email")
		if err != nil {
			http.Redirect(w, r, "/about", http.StatusFound)
			return
		}
		email := cookie.Value

		var userID int
		err = db.QueryRow(`SELECT id FROM users WHERE email=$1`, email).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			title := r.FormValue("title")
			description := r.FormValue("description")

			if title == "" {
				http.Redirect(w, r, "/tasks?error="+url.QueryEscape("Название задачи обязательно"), http.StatusFound)
				return
			}

			_, err = db.Exec(`INSERT INTO tasks (user_id, title, description) VALUES ($1, $2, $3)`,
				userID, title, description)
			if err != nil {
				http.Redirect(w, r, "/tasks?error="+url.QueryEscape("Ошибка базы данных"), http.StatusFound)
				return
			}

			http.Redirect(w, r, "/tasks", http.StatusFound)
			return
		}

		// GET - показываем список задач
		rows, err := db.Query(`SELECT id, title, description, status FROM tasks WHERE user_id=$1 ORDER BY created_at DESC`, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
			if err != nil {
				http.Error(w, "Error scanning task", http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		htmlBytes, err := os.ReadFile("tasks.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}

		html := string(htmlBytes)
		var tasksHTML string

		if len(tasks) == 0 {
			tasksHTML = `<div class="empty"><p>У вас пока нет задач</p><p style="font-size:14px;">Создайте свою первую задачу выше!</p></div>`
		} else {
			for _, task := range tasks {
				statusText := "Активно"
				statusClass := "active"
				if task.Status == "done" {
					statusText = "Выполнено"
					statusClass = "done"
				}

				titleClass := ""
				if task.Status == "done" {
					titleClass = " done"
				}

				toggleText := "Выполнить"
				if task.Status == "done" {
					toggleText = "Вернуть"
				}

				tasksHTML += `<div class="task">
					<div class="task-info">
						<div class="task-title` + titleClass + `">` + task.Title + `</div>`

				if task.Description != "" {
					tasksHTML += `<div class="task-desc">` + task.Description + `</div>`
				}

				tasksHTML += `<span class="task-status ` + statusClass + `">` + statusText + `</span>
					</div>
					<div class="task-actions">
						<form method="POST" action="/tasks/toggle">
							<input type="hidden" name="task_id" value="` + fmt.Sprint(task.ID) + `">
							<button type="submit" class="btn-toggle">` + toggleText + `</button>
						</form>
						<form method="POST" action="/tasks/delete">
							<input type="hidden" name="task_id" value="` + fmt.Sprint(task.ID) + `">
							<button type="submit" class="btn-delete">Удалить</button>
						</form>
					</div>
				</div>`
			}
		}

		// html = strings.Replace(html, "{{ if .Tasks }}", "", 1)
		// html = strings.Replace(html, "{{ else }}", "", 1)
		// html = strings.Replace(html, "{{ end }}", "", 1)
		// html = strings.Replace(html, "{{ range .Tasks }}", "", 1)
		// html = strings.Replace(html, "{{ .Title }}", "", 1)
		// html = strings.Replace(html, "{{ .Description }}", "", 1)
		// html = strings.Replace(html, "{{ .Status }}", "", 1)
		// html = strings.Replace(html, "{{ .ID }}", "", 1)

		html = strings.Replace(html, `<div id="tasks-list">`, `<div id="tasks-list">`+tasksHTML, 1)

		fmt.Fprint(w, html)
	})
	//	deletetask
	http.HandleFunc("/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("email")
		if err != nil {
			http.Redirect(w, r, "/about", http.StatusFound)
			return
		}
		email := cookie.Value

		var userID int
		err = db.QueryRow(`SELECT id FROM users WHERE email=$1`, email).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		taskID := r.FormValue("task_id")
		if taskID == "" {
			http.Redirect(w, r, "/tasks?error="+url.QueryEscape("ID задачи не указан"), http.StatusFound)
			return
		}

		_, err = db.Exec(`DELETE FROM tasks WHERE id=$1 AND user_id=$2`, taskID, userID)
		if err != nil {
			http.Redirect(w, r, "/tasks?error="+url.QueryEscape("Ошибка удаления"), http.StatusFound)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusFound)
	})
	//managetask
	http.HandleFunc("/tasks/toggle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("email")
		if err != nil {
			http.Redirect(w, r, "/about", http.StatusFound)
			return
		}
		email := cookie.Value

		var userID int
		err = db.QueryRow(`SELECT id FROM users WHERE email=$1`, email).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		taskID := r.FormValue("task_id")
		if taskID == "" {
			http.Redirect(w, r, "/tasks?error="+url.QueryEscape("ID задачи не указан"), http.StatusFound)
			return
		}

		_, err = db.Exec(`UPDATE tasks SET status = CASE 
			WHEN status = 'active' THEN 'done' 
			ELSE 'active' 
		END WHERE id=$1 AND user_id=$2`, taskID, userID)
		if err != nil {
			http.Redirect(w, r, "/tasks?error="+url.QueryEscape("Ошибка обновления статуса"), http.StatusFound)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusFound)
	})

	fmt.Println("Сервер запущен на порту http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func showAboutWithError(w http.ResponseWriter, msg string) {
	htmlBytes, err := os.ReadFile("about.html")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	html := string(htmlBytes)
	html = strings.Replace(html, `<div id="error"></div>`, `<div style="color:red">`+msg+`</div>`, 1)
	fmt.Fprint(w, html)
}

func showRegisterWithError(w http.ResponseWriter, msg string) {
	htmlBytes, err := os.ReadFile("register.html")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	html := string(htmlBytes)
	html = strings.Replace(html, `<div id="error"></div>`, `<div style="color:red">`+msg+`</div>`, 1)
	fmt.Fprint(w, html)
}
