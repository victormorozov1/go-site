package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

const (
	usersTable = "users2"
)

type User struct {
	Id                        int
	Name, Surname, Patronymic string
	Role                      string
	Phone, Email              string
	Photo_src                 string
}

func GetUsers(db *sql.DB) ([]*User, error) {
	f := func(rows *sql.Rows) (*User, error) {
		var new_user = User{}

		err := rows.Scan(&new_user.Id, &new_user.Name, &new_user.Surname, &new_user.Patronymic,
			&new_user.Role, &new_user.Phone, &new_user.Email, &new_user.Photo_src)

		return &new_user, err
	}
	return getFromDB(db, "select * from "+usersTable, f)
}

func (user *User) SaveToDB(db *sql.DB) error {
	return SaveToDB(db, "users2", map[string]string{
		"id":         strconv.Itoa(user.Id),
		"name":       user.Name,
		"surname":    user.Surname,
		"patronymic": user.Patronymic,
		"role":       user.Role,
		"phone":      user.Phone,
		"email":      user.Email,
		"photo_src":  user.Photo_src,
	})
}

func (user *User) String() string {
	return fmt.Sprintf("User(id=%d name=%v surname=%v patronymic=%v role=%v phone=%v email=%v photo_src=%v",
		user.Id, user.Name, user.Surname, user.Patronymic, user.Role, user.Phone, user.Email, user.Photo_src)
}
