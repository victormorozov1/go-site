package server

import "fmt"

type Session struct {
	Id   int
	Name string
}

func (session Session) String() string {
	return fmt.Sprintf("Session#%d(name: %s)", session.Id, session.Name)
}
