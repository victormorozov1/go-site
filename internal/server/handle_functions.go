package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
)

func Hash(s string) string {
	return s + "типо хэширую" // Тут по нормальному хеширование сделать
}

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{}, r))
	if err != nil {
		println(err.Error())
	}
}

func (server *Server) TestPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/test.html")
	if err != nil {
		print(err)
	}
	println("returning in test", server.BaseTemplateData)
	err = t.Execute(w, server.BaseTemplateData)
	if err != nil {
		println(err.Error())
	}
}

func (server *Server) MapPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/map.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		print(err)
	}
	println("returning in map", server.BaseTemplateData)
	err = t.Execute(w, server.BaseTemplateData)
	if err != nil {
		println(err.Error())
	}
}

func (server *Server) WorkplacePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/workplace.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		print(err)
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["place_id"])
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	workplace, err := database.GetTableById(&server.DataBase, id)
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	data := server.GetTemplateAndUserData([]*map[string]interface{}{{"Workplace": workplace}}, r)
	println("returning in workplace", data)
	err = t.Execute(w, data)
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}
}
