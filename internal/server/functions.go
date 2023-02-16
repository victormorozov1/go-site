package server

import (
	"errors"
	"internal/db"
	"net/http"
	"strconv"
)

func UsernameExists(server *Server, name string) (bool, error) {
	allUsers, err := db.GetAllUsers(&server.DataBase)
	if err != nil {
		return false, err
	}

	for _, user := range allUsers {
		if user.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func RegisterUserCheck(server *Server, user *db.User) error {
	err := user.Check()
	if err != nil {
		return err
	}

	usernameExists, err := UsernameExists(server, user.Name)
	if err != nil {
		panic(err.Error())
		return err
	}

	if usernameExists {
		return errors.New("Username already exists")
	}

	return nil
}

func GetUserId(server *Server, r *http.Request) (int, error) {
	cookie, err := r.Cookie(server.CookieName)
	if err != nil {
		return -1, err
	}

	var sessionId int
	sessionId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return -1, err
	}

	session, ok := server.Sessions[sessionId]
	if !ok {
		return -1, errors.New("session#" + strconv.Itoa(sessionId) + " not found")
	}

	return session.UserId, nil
}
