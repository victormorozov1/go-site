package server

import (
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
)

func (server *Server) ReservationPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/reservation.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		print(err)
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	reservation, err := database.GetReservationById(&server.DataBase, id)
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	err = reservation.LoadUser(&server.DataBase)
	if err != nil {
		println(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	err = reservation.LoadTable(&server.DataBase)
	if err != nil {
		println(err.Error())
		t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
		return
	}

	err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Reservation": reservation}}, r))
	if err != nil {
		println(err.Error())
	}
}

func (server *Server) DeleteReservationAjaxHandler(w http.ResponseWriter, r *http.Request) {
	// Тут нужно по нормальному возвращать json а не строчки
	errorFunc := func(err error) {
		println(err.Error())
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			println(err.Error())
		}
	}

	if r.Method == "POST" {
		userId, err := GetUserId(server, r)
		if err != nil {
			errorFunc(err)
			return
		}

		reservationId, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			errorFunc(err)
			return
		}

		println("Request to delete reservation #" + strconv.Itoa(reservationId))

		reservation, err := database.GetReservationById(&server.DataBase, reservationId)
		if err != nil {
			errorFunc(err)
			return
		}

		user, err := database.GetUserById(&server.DataBase, userId)
		if err != nil {
			errorFunc(err)
			return
		}
		if user.Id != reservation.User_id && user.Role != server.Roles.Admin {
			errorFunc(errors.New("Недостаточно прав"))
			return
		}

		reservation.Delete(&server.DataBase) // ЛУчше сделать просто функцию deleteReservationById

		_, err = w.Write([]byte("ok"))
		if err != nil {
			println(err.Error())
		}
	}
}
