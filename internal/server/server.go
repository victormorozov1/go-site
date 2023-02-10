package server

import (
	"database/sql"
	"fmt"
	"golang.org/x/exp/maps"
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
}

func (server *Server) handleFunc() {
	http.HandleFunc(server.Routes.MainPage, server.mainPage)
	http.HandleFunc(server.Routes.AllUsersPage, server.allUsersPage)
	http.HandleFunc(server.Routes.RegisterPage, server.Register)
	http.HandleFunc(server.Routes.LoginPage, server.LogIn)
	http.HandleFunc(server.Routes.UserCabinet, server.UserPage)
	http.HandleFunc(server.Routes.TestPage, server.TestPage)
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

func AddRoutesData(data *map[string]interface{}, server *Server) {
	(*data)["Routes"] = server.Routes
}

func JoinData(data1 *map[string]interface{}, data2 *map[string]interface{}) *map[string]interface{} { // Вынести куда-нибудь
	resData := make(map[string]interface{})
	maps.Copy(resData, *data1)
	maps.Copy(resData, *data2)
	return &resData
}
