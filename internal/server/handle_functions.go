package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
	"time"
)

func Hash(s string) string {
	return s + "типо хэширую" // Тут по нормальному хеширование сделать
}

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.Execute(w, server.BaseTemplateData)
}

func (server *Server) allUsersPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/all_users.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
		return
	}

	allUsers, err := database.GetAllUsers(server.DataBase, server.UsersTableName, server.ReservationsTableName)

	if err != nil {
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
		println("Error: " + err.Error())
		return
	}

	t.Execute(w, server.GetTemplateData(&map[string]interface{}{
		"UsersNum": len(allUsers),
		"UsersArr": allUsers,
	})) // Тут нельзя возвращать пароли
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/register.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		http.Redirect(w, r, server.Routes.MainPage, http.StatusSeeOther)
		panic(err)
		return
	}

	if r.Method == "POST" {
		name, password, repeatPassword := r.FormValue("name"), r.FormValue("password"), r.FormValue("repeat-password")
		if password == repeatPassword {
			newUser := database.User{
				Name:           name,
				HashedPassword: Hash(password),
			}

			err = RegisterUserCheck(server, &newUser)
			if err != nil {
				print(err)
				t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
				return
			}

			println(newUser.String())
			err = newUser.SaveToDB(server.DataBase, server.UsersTableName)
			if err != nil {
				println(err)
			} else {
				http.Redirect(w, r, server.Routes.AllUsersPage, http.StatusSeeOther)
			}
		} else {
			t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": "Passwords don't match"}))
		}
	} else {
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": ""}))
	}
}

func (server *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		println(err.Error())
		http.Redirect(w, r, server.Routes.MainPage, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		println("In login")
		name, password := r.FormValue("name"), r.FormValue("password")
		dbUsers, err := database.GetUserBy(server.DataBase, "name", name, server.UsersTableName)

		if err != nil {
			t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": "Error on server"}))
			panic(err) // Хз что тут может быть, но можно как-то обработать
			return
		}

		if len(dbUsers) == 0 {
			t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": "User not found"}))
			println("Users with name " + name + " not found")
			return
		}

		if len(dbUsers) > 1 {
			panic("Error: found several users with name " + name) // тоже нудно вернуть ошибку пользователю
		}

		dbUser := dbUsers[0]

		if Hash(password) == dbUser.HashedPassword {
			println(dbUser.Name + " logged in successfully")

			session := server.CreateSession()
			session.UserId = dbUser.Id
			println("Created session id = " + session.String())

			cookie := &http.Cookie{
				Name:    server.CookieName,
				Value:   strconv.Itoa(session.Id),
				Expires: time.Now().Add(time.Minute * 10),
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, server.Routes.UserCabinet, http.StatusSeeOther)

			return
		} else {
			println("wrong password")
			t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": "Wrong password"}))
			return
		}
	} else {
		t.Execute(w, server.GetTemplateData())
		return
	}
}

func (server *Server) UserPage(w http.ResponseWriter, r *http.Request) {
	var err error
	var t *template.Template

	t, err = template.ParseFiles("templates/user.html", "templates/navbar.html", "templates/include.html")

	defer func() { // Наверно это убрать лучше
		if err != nil {
			println(err)
		}
	}()

	var cookie *http.Cookie
	cookie, err = r.Cookie(server.CookieName)
	if err != nil {
		t.Execute(w, err.Error())
		// Вернуть ошибку
		return
	}

	var sessionId int
	sessionId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		t.Execute(w, err.Error())
		// Вернуть ошибку
		return
	}

	var (
		session *Session
		ok      bool
	)
	session, ok = server.Sessions[sessionId]
	if !ok {
		println("Session#" + strconv.Itoa(sessionId) + "not found")
		// Вернуть ошибку
		return
	}

	var user *database.User
	user, err = database.GetUserById(server.DataBase, session.UserId, server.UsersTableName)

	err = user.LoadReservationsFromDB(server.DataBase, server.ReservationsTableName)
	if err != nil {
		print(err.Error())
	}

	data := map[string]interface{}{
		"Phone":        user.Phone,
		"Email":        user.Email,
		"PhotoSrc":     user.Photo_src,
		"Role":         user.Role,
		"Patronymic":   user.Patronymic,
		"Name":         user.Name,
		"Surname":      user.Surname,
		"Reservations": user.Reservations,
	}

	t.Execute(w, server.GetTemplateData(&data))
}

func (server *Server) TestPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/test.html")
	if err != nil {
		print(err)
	}
	println("returning in test", server.BaseTemplateData)
	t.Execute(w, server.BaseTemplateData)
}

func (server *Server) ReservationPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/reservation.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		print(err)
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
		return
	}

	reservation, err := database.GetReservationById(server.DataBase, server.ReservationsTableName, id)
	if err != nil {
		print(err.Error())
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
		return
	}

	err = reservation.LoadUser(server.DataBase, server.UsersTableName)
	if err != nil {
		println(err.Error())
		t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Error": err.Error()}))
		return
	}

	t.Execute(w, server.GetTemplateData(&map[string]interface{}{"Reservation": reservation}))
}
