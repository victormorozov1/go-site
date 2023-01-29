package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Host string
	Port int
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
}

func (server *Server) Start() {
	fmt.Print("Start ", server)
	server.handleFunc()
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}
