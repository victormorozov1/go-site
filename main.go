package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Server struct {
	host string
	port int
}

func (server *Server) String() string {
	return fmt.Sprintf("Server(%s:%d)", server.host, server.port)
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

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.Execute(w, nil)
}

func (server *Server) handleFunc() {
	http.HandleFunc("/", server.mainPage)
}

func (server *Server) start() {
	fmt.Print("Start ", server)
	//http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)
	http.ListenAndServe(":8080", nil)
}

func main() {
	s := Server{"127.0.0.1", 8080}
	s.handleFunc()
	s.start()
}
