package server

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
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
	Routes                                *Routes
	BaseTemplateData                      map[string]interface{}
}

type Routes struct {
	MainPage                                           string
	AllUsersPage, UserCabinet, RegisterPage, LoginPage string
	TestPage                                           string
	ReservationPage                                    string
}

func (server *Server) handleFunc() {
	http.HandleFunc(server.Routes.MainPage, server.mainPage)
	http.HandleFunc(server.Routes.AllUsersPage, server.allUsersPage)
	http.HandleFunc(server.Routes.RegisterPage, server.Register)
	http.HandleFunc(server.Routes.LoginPage, server.LogIn)
	http.HandleFunc(server.Routes.UserCabinet, server.UserPage)
	http.HandleFunc(server.Routes.TestPage, server.TestPage)

	rtr := mux.NewRouter()
	rtr.HandleFunc(server.Routes.ReservationPage+"/{id:[0-9]+}", server.ReservationPage).Methods("GET")

	http.Handle("/", rtr)
}

func (server *Server) String() string {
	return fmt.Sprintf("Server(%s:%d)", server.Host, server.Port)
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
