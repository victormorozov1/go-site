package server

import (
	"fmt"
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
	"time"
)

func Hash(s string) string {
	return s + "типо хэширую" // Тут по нормальному хеширование сделать
}

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.Execute(w, nil)
}

func (server *Server) allUsersPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/all_users.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	allUsers, err := database.GetAllUsers(server.DataBase, server.UsersTableName, server.ReservationsTableName)

	if err != nil {
		panic(err) // Тут нужно написать возвращение ошибки пользователю
	}

	type templateData struct {
		UsersNum int
		UsersArr []*database.User
	}
	t.Execute(w, &templateData{len(allUsers), allUsers}) // Тут нельзя возвращать пароли
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name, password, repeatPassword := r.FormValue("name"), r.FormValue("password"), r.FormValue("repeat-password")
		if password == repeatPassword {
			newUser := database.User{
				Name:           name,
				HashedPassword: Hash(password),
			}
			println(newUser.String())
			err := newUser.SaveToDB(server.DataBase, server.UsersTableName)
			if err != nil {
				println(err)
				// Нужно вернуть ошибку
			} else {
				// редирект регистрация успешна
			}
		} else {
			// Нужно вернуть ошибку ПАРОЛИ НЕ СОВПАДАЮТ
		}
	} else {
		t, err := template.ParseFiles("templates/register.html")

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.Execute(w, nil)
	}
}

func (server *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		println("In login")
		name, password := r.FormValue("name"), r.FormValue("password")
		dbUsers, err := database.GetUserBy(server.DataBase, "name", name,
			server.UsersTableName, server.ReservationsTableName)

		if err != nil {
			panic(err) // Хз что тут может быть, но можно как-то обработать
		}

		if len(dbUsers) == 0 {
			// Нужно вернуть ошибку, что пользователь не найден
			println("Users with name " + name + " not found")
			return
		}

		if len(dbUsers) > 1 {
			panic("Error: found several users with name " + name) // тоже нудно вернуть ошибку пользователю
		}

		dbUser := dbUsers[0]

		if Hash(password) == dbUser.HashedPassword {
			println(dbUser.Name + " logged in successfully")

			// Можно добавить структуру сессии, если много данных нужно будет хранить
			session := server.CreateSession()
			session.Name = name
			println("Created session id = " + session.String())

			cookie := &http.Cookie{
				Name:    server.CookieName,
				Value:   strconv.Itoa(session.Id), // Нужно хотя-бы имя хешировать, присем не так как пароль
				Expires: time.Now().Add(time.Minute * 10),
			}

			http.SetCookie(w, cookie)
			// чото тоже вернуть надо
		} else {
			println("wrong password")
			// Нужно вернуть ошибку пароль неверный
		}
	} else {
		t, err := template.ParseFiles("templates/login.html")

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.Execute(w, nil)
	}
}

//func UserPage(w http.ResponseWriter, r *http.Request) {
//
//}
