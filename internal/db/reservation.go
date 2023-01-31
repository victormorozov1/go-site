package db

type Reservation struct {
	id, table_id         int
	start_time, end_time int // Потом сделать Time
}
