package server

import "fmt"

type Session struct {
	Id     int
	UserId int
}

func (session Session) String() string {
	return fmt.Sprintf("Session#%d(user_id: %v)", session.Id, session.UserId)
}
