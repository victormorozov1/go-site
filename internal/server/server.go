package server

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"internal/db"
	"math/rand"
	"net/http"
)

type Server struct {
	Host                                  string
	Port                                  int
	DataBase                              db.Database
	UsersTableName, ReservationsTableName string
	CookieName                            string
	Sessions                              map[int]*Session
	Routes                                *Routes
	BaseTemplateData                      *map[string]interface{}
	Roles                                 Roles
}

type Routes struct {
	MainPage                                           string
	AllUsersPage, UserCabinet, RegisterPage, LoginPage string
	TestPage                                           string
	ReservationPage, DeleteReservationAjaxHandler      string
	MapPage, WorkplacePage                             string
}

type Roles struct {
	User, Admin string
}

func (routes Routes) String() string {
	return fmt.Sprintf("Routes(MainPage: %s, AllUsersPage: %s, UserCabinet: %s, RegisterPage: %s, LoginPage: %s, TestPage: %s, ReservationPage: %s",
		routes.MainPage, routes.AllUsersPage, routes.UserCabinet, routes.RegisterPage, routes.LoginPage, routes.TestPage, routes.ReservationPage)
}

func (server *Server) handleFunc() {
	http.HandleFunc(server.Routes.MainPage, server.mainPage)
	http.HandleFunc(server.Routes.AllUsersPage, server.AllUsersPage)
	http.HandleFunc(server.Routes.RegisterPage, server.Register)
	http.HandleFunc(server.Routes.LoginPage, server.LogIn)
	http.HandleFunc(server.Routes.UserCabinet, server.UserPage)
	http.HandleFunc(server.Routes.TestPage, server.TestPage)
	http.HandleFunc(server.Routes.DeleteReservationAjaxHandler, server.DeleteReservationAjaxHandler)
	http.HandleFunc(server.Routes.MapPage, server.MapPage)

	rtr := mux.NewRouter()
	rtr.HandleFunc(server.Routes.ReservationPage+"/{id:[0-9]+}", server.ReservationPage).Methods("GET")
	rtr.HandleFunc(server.Routes.WorkplacePage+"/{place_id:[0-9]+}", server.WorkplacePage).Methods("GET")

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
	server.DataBase.Connection, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/go_site")
	if err != nil {
		panic(err)
	}
	defer server.DataBase.Connection.Close()

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
