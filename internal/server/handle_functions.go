package server

import (
	"fmt"
	"html/template"
	database "internal/db"
	"net/http"
)

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

	allUsers, err := database.GetUsers(server.DataBase)

	if err != nil {
		panic(err) // Тут нужно написать возвращение ошибки пользователю
	}

	type templateData struct {
		UsersNum int
	}
	t.Execute(w, &templateData{len(allUsers)})
}
