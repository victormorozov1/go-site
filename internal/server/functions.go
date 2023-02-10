package server

import (
	"errors"
	"internal/db"
)

func UsernameExists(server *Server, name string) (bool, error) {
	allUsers, err := db.GetAllUsers(server.DataBase, server.UsersTableName, server.ReservationsTableName)
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
