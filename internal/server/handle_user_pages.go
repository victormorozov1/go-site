package server

import (
	"html/template"
	database "internal/db"
	"net/http"
	"strconv"
	"time"
)

func (server *Server) AllUsersPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/all_users.html", "templates/navbar.html", "templates/include.html")

	if err != nil {
		println("Error: ", err.Error())
		err = t.Execute(w, server.GetTemplateAndUserData(
			[]*map[string]interface{}{
				{"Error": err.Error()}},
			r))
		if err != nil {
			println(err.Error())
		}
		return
	}

	allUsers, err := database.GetAllUsers(&server.DataBase)

	if err != nil {
		println("Error: " + err.Error())
		err = t.Execute(w, server.GetTemplateAndUserData(
			[]*map[string]interface{}{
				{"Error": err.Error()}},
			r))
		if err != nil {
			println(err.Error())
		}
		return
	}

	m := &map[string]interface{}{
		"UsersNum": len(allUsers),
		"UsersArr": allUsers,
	}
	err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{m}, r))
	if err != nil {
		println(err.Error())
	}
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
				err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": err.Error()}}, r))
				if err != nil {
					println(err)
				}
				return
			}

			println(newUser.String())
			err = newUser.SaveToDB(&server.DataBase)
			if err != nil {
				println(err)
			} else {
				http.Redirect(w, r, server.Routes.AllUsersPage, http.StatusSeeOther)
			}
		} else {
			err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": "Passwords don't match"}}, r))
			if err != nil {
				println(err)
			}
		}
	} else {
		err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": ""}}, r))
		if err != nil {
			print(err)
		}
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
		dbUsers, err := database.GetUserBy(&server.DataBase, server.DataBase.Tables.Users.Name, name)

		if err != nil {
			println(err)
			err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": "Error on server"}}, r))
			if err != nil {
				println(err)
			}
			return
		}

		if len(dbUsers) == 0 {
			println("Users with name " + name + " not found")
			err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": "User not found"}}, r))
			if err != nil {
				println(err)
			}
			return
		}

		if len(dbUsers) > 1 {
			println("Error: found several users with name " + name) // тоже нудно вернуть ошибку пользователю
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
			err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{{"Error": "Wrong password"}}, r))
			if err != nil {
				println(err)
			}
			return
		}
	} else {
		err = t.Execute(w, server.GetTemplateAndUserData(nil, r))
		if err != nil {
			println(err)
		}
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
		http.Redirect(w, r, server.Routes.LoginPage, http.StatusSeeOther)
		return
	}

	var sessionId int
	sessionId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		println(err)
		err = t.Execute(w, err.Error())
		// Вернуть ошибку
		if err != nil {
			println(err)
		}
		return
	}

	var (
		session *Session
		ok      bool
	)
	session, ok = server.Sessions[sessionId]
	if !ok {
		http.Redirect(w, r, server.Routes.LoginPage, http.StatusSeeOther)
		return
	}

	var user *database.User
	user, err = database.GetUserById(&server.DataBase, session.UserId)

	err = user.LoadReservationsFromDB(&server.DataBase)
	if err != nil {
		print(err.Error())
	}

	err = t.Execute(w, server.GetTemplateAndUserData([]*map[string]interface{}{}, r))
	if err != nil {
		println(err)
	}
}
