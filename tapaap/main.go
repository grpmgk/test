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

	_ "github.com/lib/pq"
)

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
