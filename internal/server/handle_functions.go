package server

import (
	"fmt"
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
	"time"
)

func Hash(s string) string {
	return s + "типо хэширую" // Тут по нормальному хеширование сделать
}

type errorStruct struct {
	Error string
}

func (server *Server) mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.Execute(w, nil)
}

func (server *Server) allUsersPage(w http.ResponseWriter, r *http.Request) {
	type templateData struct {
		UsersNum int
		UsersArr []*database.User
		Error    string
	}

	t, err := template.ParseFiles("templates/all_users.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		t.Execute(w, templateData{Error: err.Error()})
		return
	}

	allUsers, err := database.GetAllUsers(server.DataBase, server.UsersTableName, server.ReservationsTableName)

	if err != nil {
		t.Execute(w, templateData{Error: err.Error()})
		println("Error: " + err.Error())
		return
	}

	t.Execute(w, &templateData{
		UsersNum: len(allUsers),
		UsersArr: allUsers,
	}) // Тут нельзя возвращать пароли
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/register.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

			err = newUser.Check()
			if err != nil {
				print(err)
				t.Execute(w, errorStruct{err.Error()})
				return
			}

			println(newUser.String())
			err = newUser.SaveToDB(server.DataBase, server.UsersTableName)
			if err != nil {
				println(err)

			} else {
				http.Redirect(w, r, "/users", http.StatusSeeOther)
			}
		} else {
			t.Execute(w, errorStruct{"Passwords don't match"})
		}
	} else {
		t.Execute(w, errorStruct{""})
	}
}

func (server *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/navbar.html", "templates/include.html")
	if err != nil {
		println(err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		println("In login")
		name, password := r.FormValue("name"), r.FormValue("password")
		dbUsers, err := database.GetUserBy(server.DataBase, "name", name,
			server.UsersTableName, server.ReservationsTableName)

		if err != nil {
			t.Execute(w, errorStruct{"Error on server"})
			panic(err) // Хз что тут может быть, но можно как-то обработать
			return
		}

		if len(dbUsers) == 0 {
			t.Execute(w, errorStruct{"User not found"})
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
			http.Redirect(w, r, "/me", http.StatusSeeOther)
		} else {
			println("wrong password")
			t.Execute(w, errorStruct{"Wrong password"})
		}
	} else {
		t.Execute(w, nil)
	}
}

func (server *Server) UserPage(w http.ResponseWriter, r *http.Request) {
	var err error

	var t *template.Template
	t, err = template.ParseFiles("templates/user.html", "templates/navbar.html", "templates/include.html")

	type TemplateData struct {
		Name, Surname, Patronymic string
		Role                      string
		Phone, Email              string
		PhotoSc                   string
	}
	var templateData TemplateData

	defer func() {
		t.Execute(w, templateData)
		if err != nil {
			println(err)
		}
	}()

	var cookie *http.Cookie
	cookie, err = r.Cookie(server.CookieName)
	if err != nil {
		return
	}

	var sessionId int
	sessionId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return
	}

	var (
		session *Session
		ok      bool
	)
	session, ok = server.Sessions[sessionId]
	if !ok {
		println("Session#" + strconv.Itoa(sessionId) + "not found")
		return
	}

	var users []*database.User
	users, err = database.GetUserBy(server.DataBase, "id", strconv.Itoa(session.UserId),
		server.UsersTableName, server.ReservationsTableName) // добавить функцию GetUserById
	user := users[0]

	templateData.Phone = user.Phone
	templateData.Email = user.Email
	templateData.PhotoSc = user.Photo_src
	templateData.Role = user.Role
	templateData.Patronymic = user.Patronymic
	templateData.Name = user.Name
	templateData.Surname = user.Surname

	print("returning template data .name " + user.Surname)
}
