package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id                        int
	Name, Surname, Patronymic string
	Role                      string
	Phone, Email              string
	Photo_src                 string
	Reservations              []Reservation
}

//func GetUsers(db *sql.DB) (error, *[]*User) {
//	result, err := db.Query("select * from users")
//
//	var all_users []*User
//
//	if err != nil {
//		return err, &all_users
//	}
//
//	for result.Next() {
//		var new_user = User{}
//		err := result.Scan(&new_user.id, &new_user.role, &new_user.name, &new_user.surname, &new_user.patronymic,
//			&new_user.phone, &new_user.email, &new_user.photo_src)
//
//		if err != nil {
//			return err, &all_users
//		}
//
//		all_users = append(all_users, &new_user)
//	}
//
//	return nil, &all_users
//}

func (user *User) SaveToDB(db *sql.DB) error {
	insert, err := db.Query("INSERT INTO `users` (`id`, `name`, `surname`, `patronymic`, `role`, `phone`, `email`, `photo_src`)" +
		fmt.Sprintf(" VALUES('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')",
			user.Id, user.Name, user.Surname, user.Patronymic, user.Role, user.Phone, user.Email, user.Photo_src)) // add reservations
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}
