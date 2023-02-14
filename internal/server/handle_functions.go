package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func Hash(s string) string {
	return s + "типо хэширую" // Тут по нормальному хеширование сделать
}

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{}, r))
}

func (server *Server) TestPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/test.html")
	if err != nil {
		print(err)
	}
	println("returning in test", server.BaseTemplateData)
	t.Execute(w, server.BaseTemplateData)
}
