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

func GetUsers(db *sql.DB, usersTableName, reservationsTableName, criteria string) ([]*User, error) {
	f := func(rows *sql.Rows) (*User, error) {
		return ScanUserAndReservationsFromDBRows(rows, db, reservationsTableName)
	}
	return getFromDB(db, "select * from "+usersTableName+" WHERE "+criteria, f)
}

func GetAllUsers(db *sql.DB, usersTableName, reservationsTableName string) ([]*User, error) {
	return GetUsers(db, usersTableName, reservationsTableName, "")
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

func (user *User) LoadReservationsFromDB(db *sql.DB, reservationsTableName string) error {
	reservations, err := GetReservationsByUserId(db, reservationsTableName, user.Id)

	if err != nil {
		return err
	}

	user.Reservations = reservations

	return nil
}

func ScanUserFromDBRows(rows *sql.Rows) (*User, error) {
	var newUser = User{}

	err := rows.Scan(&newUser.Id, &newUser.Name, &newUser.Surname, &newUser.Patronymic,
		&newUser.Role, &newUser.Phone, &newUser.Email, &newUser.Photo_src, &newUser.HashedPassword)

	return &newUser, err
}

func ScanUserAndReservationsFromDBRows(rows *sql.Rows, db *sql.DB, reservationsTableName string) (*User, error) {
	newUser, err := ScanUserFromDBRows(rows)

	if err != nil {
		return nil, err
	}

	newUser.LoadReservationsFromDB(db, reservationsTableName)

	return newUser, err
}

func GetUserBy(db *sql.DB, columnName, value, usersTableName, reservationsTableName string) ([]*User, error) {
	return GetUsers(db, usersTableName, reservationsTableName, columnName+"='"+value+"'")
}
