package db_structs

type user struct {
	id                        int
	name, surname, patronymic string
	role                      string
	phone, email              string
	photo_src                 string
	reservations              []Reservation
}
