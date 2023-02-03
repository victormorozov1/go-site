package db

import (
	"database/sql"
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
	Reservations              []Reservation
}

//	func GetUsers(db *sql.DB) (error, *[]*User) {
//		result, err := db.Query("select * from users")
//
//		var all_users []*User
//
//		if err != nil {
//			return err, &all_users
//		}
//
//		for result.Next() {
//			var new_user = User{}
//			err := result.Scan(&new_user.id, &new_user.role, &new_user.name, &new_user.surname, &new_user.patronymic,
//				&new_user.phone, &new_user.email, &new_user.photo_src)
//
//			if err != nil {
//				return err, &all_users
//			}
//
//			all_users = append(all_users, &new_user)
//		}
//
//		return nil, &all_users
//	}
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
