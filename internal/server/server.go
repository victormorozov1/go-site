package server

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Server struct {
	Host                                  string
	Port                                  int
	DataBase                              *sql.DB
	UsersTableName, ReservationsTableName string
}

func (server *Server) String() string {
	return fmt.Sprintf("Server(%s:%d)", server.Host, server.Port)
}

//func (server *Server) parseFiles(w http.ResponseWriter, fileNames ...string) *template.Template {
//	t, err := template.ParseFiles(fileNames...)
//
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//	}
//
//	return t
//}

func (server *Server) handleFunc() {
	http.HandleFunc("/", server.mainPage)
	http.HandleFunc("/users", server.allUsersPage)
	http.HandleFunc("/register", server.Register)
	http.HandleFunc("/login", server.LogIn)
}

func (server *Server) Start() {
	fmt.Print("Start ", server)
	server.handleFunc()

	var err error
	server.DataBase, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_site")
	if err != nil {
		panic(err)
	}
	defer server.DataBase.Close()

	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}
