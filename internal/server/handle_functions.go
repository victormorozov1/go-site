package server

import (
	"fmt"
	"html/template"
	database "internal/db"
	"net/http"
)

func Hash(s string) string {
	return s // Тут по нормальному хеширование сделать
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

	allUsers, err := database.GetUsers(server.DataBase, server.UsersTableName, server.ReservationsTableName)

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
			err := newUser.SaveToDB(server.DataBase)
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
