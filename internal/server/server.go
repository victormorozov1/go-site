package server

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
)

type Server struct {
	Host                                  string
	Port                                  int
	DataBase                              *sql.DB
	UsersTableName, ReservationsTableName string
	CookieName                            string
	Sessions                              map[int]*Session
}

func (server *Server) String() string {
	return fmt.Sprintf("Server(%s:%d)", server.Host, server.Port)
}

func (server *Server) handleFunc() {
	http.HandleFunc("/", server.mainPage)
	http.HandleFunc("/users", server.allUsersPage)
	http.HandleFunc("/register", server.Register)
	http.HandleFunc("/login", server.LogIn)
}

func (server *Server) Start() {
	if server.Sessions == nil {
		server.Sessions = make(map[int]*Session)
	}

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

func (server *Server) CreateSession() *Session {
	id := 0
	ok := true
	for ok {
		id = rand.Int()
		_, ok = server.Sessions[id]
	}
	server.Sessions[id] = &Session{Id: id}
	return server.Sessions[id]
}
