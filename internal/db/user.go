package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type User struct {
	Id                        int
	Name, Surname, Patronymic string
	Role                      string
	Phone, Email              string
	Photo_src                 string
	HashedPassword            string
	Reservations              []*Reservation
}

func GetUsers(db *sql.DB, usersTableName, reservationsTableName string) ([]*User, error) {
	f := func(rows *sql.Rows) (*User, error) {
		var newUser = User{}

		err := rows.Scan(&newUser.Id, &newUser.Name, &newUser.Surname, &newUser.Patronymic,
			&newUser.Role, &newUser.Phone, &newUser.Email, &newUser.Photo_src)

		if err != nil {
			return nil, err
		}

		reservations, err := GetReservationsByUserId(db, reservationsTableName, newUser.Id)

		if err != nil {
			return &newUser, err
		}

		newUser.Reservations = reservations

		return &newUser, nil
	}
	return getFromDB(db, "select * from "+usersTableName, f)
}

func (user *User) SaveToDB(db *sql.DB, usersTableName string) error {
	return SaveToDB(db, usersTableName, map[string]string{
		"id":              strconv.Itoa(user.Id), //По идее если id=0 то его не нужно отправлять, но оно и так игнорируется почему-то
		"name":            user.Name,
		"surname":         user.Surname,
		"patronymic":      user.Patronymic,
		"role":            user.Role,
		"phone":           user.Phone,
		"email":           user.Email,
		"photo_src":       user.Photo_src,
		"hashed_password": user.HashedPassword,
	})
}

func (user *User) String() string {
	return fmt.Sprintf("User(id=%d name=%v surname=%v patronymic=%v role=%v phone=%v email=%v photo_src=%v %d reservations",
		user.Id, user.Name, user.Surname, user.Patronymic, user.Role, user.Phone, user.Email, user.Photo_src, len(user.Reservations))
}
