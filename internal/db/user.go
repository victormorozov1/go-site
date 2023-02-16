package db

import (
	"database/sql"
	"errors"
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

func (user *User) Check() error { // Можно сделать получше, я пока набросал
	if user.Name == "" {
		return errors.New("Empty name")
	}

	return nil
}

func GetUsers(db *Database, criteria string) ([]*User, error) {
	f := func(rows *sql.Rows) (*User, error) {
		return ScanUserFromDBRows(rows)
	}
	if criteria != "" {
		criteria = " WHERE " + criteria
	}
	return getFromDB(db.Connection, "select * from "+db.Tables.Users.TableName+criteria, f)
}

func GetAllUsers(db *Database) ([]*User, error) {
	return GetUsers(db, "")
}

func (user *User) SaveToDB(db *Database) error {
	return SaveToDB(db.Connection, db.Tables.Users.TableName, map[string]string{
		db.Tables.Users.Id:             strconv.Itoa(user.Id), //По идее если id=0 то его не нужно отправлять, но оно и так игнорируется почему-то
		db.Tables.Users.Name:           user.Name,
		db.Tables.Users.Surname:        user.Surname,
		db.Tables.Users.Patronymic:     user.Patronymic,
		db.Tables.Users.Role:           user.Role,
		db.Tables.Users.Phone:          user.Phone,
		db.Tables.Users.Email:          user.Email,
		db.Tables.Users.PhotoSrc:       user.Photo_src,
		db.Tables.Users.HashedPassword: user.HashedPassword,
	})
}

func (user *User) String() string {
	return fmt.Sprintf("User(Id=%d Name=%v Surname=%v Patronymic=%v Role=%v Phone=%v Email=%v Photo_src=%v Reservations=%v)",
		user.Id, user.Name, user.Surname, user.Patronymic, user.Role, user.Phone, user.Email, user.Photo_src, user.Reservations)
}

func (user *User) LoadReservationsFromDB(db *Database) error {
	reservations, err := GetReservationsByUserId(db, user.Id)

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

func ScanUserAndReservationsFromDBRows(rows *sql.Rows) (*User, error) {
	newUser, err := ScanUserFromDBRows(rows)

	if err != nil {
		return nil, err
	}

	return newUser, err
}

func GetUserBy(db *Database, columnName, value string) ([]*User, error) {
	return GetUsers(db, columnName+"='"+value+"'")
}

func GetUserById(db *Database, id int) (*User, error) {
	users, err := GetUserBy(db, "id", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("User not found")
	}

	if len(users) == 0 {
		err := errors.New("many users with id=" + strconv.Itoa(id))
		panic(err.Error())
		return nil, err
	}

	return users[0], nil
}
