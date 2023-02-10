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

type Error struct { // Вынести в отдельный файл
	Message string
}

func (error *Error) Error() string {
	return error.Message
}

func (user *User) Check() error { // Можно сделать получше, я пока набросал
	if user.Name == "" {
		return &Error{"Empty name"}
	}

	return nil
}

func GetUsers(db *sql.DB, usersTableName, reservationsTableName, criteria string) ([]*User, error) {
	f := func(rows *sql.Rows) (*User, error) {
		return ScanUserAndReservationsFromDBRows(rows, db, reservationsTableName)
	}
	if criteria != "" {
		criteria = " WHERE " + criteria
	}
	return getFromDB(db, "select * from "+usersTableName+criteria, f)
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

func GetUserById(db *sql.DB, id int, usersTableName, reservationsTableName string) (*User, error) {
	users, err := GetUserBy(db, "id", strconv.Itoa(id), usersTableName, reservationsTableName)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, &Error{"User not found"}
	}

	if len(users) == 0 {
		err := Error{"many users with id=" + strconv.Itoa(id)}
		panic(err.Error())
		return nil, &err
	}

	return users[0], nil
}
