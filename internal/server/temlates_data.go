package server

import (
	"golang.org/x/exp/maps"
	"internal/db"
	"net/http"
	"strconv"
)

func JoinDataArr(dataToJoin []*map[string]interface{}) *map[string]interface{} { // Вынести куда-нибудь
	resData := make(map[string]interface{})

	for _, data := range dataToJoin {
		maps.Copy(resData, *data)
	}

	return &resData
}

func JoinData(dataToJoin ...*map[string]interface{}) *map[string]interface{} { // Вынести куда-нибудь
	return JoinDataArr(dataToJoin)
}

func (server *Server) CountBaseTemplateData() {
	server.BaseTemplateData = &map[string]interface{}{"Routes": server.Routes}
}

func (server *Server) GetTemplateData(dataToAdd ...*map[string]interface{}) *map[string]interface{} {
	joinedDataToAdd := JoinDataArr(dataToAdd)
	return JoinData(joinedDataToAdd, server.BaseTemplateData)
}

func UserAuthorizedData(val bool) *map[string]interface{} {
	return &map[string]interface{}{"Authorized": val}
}

func GetTemplateUserData(server *Server, r *http.Request) *map[string]interface{} {
	cookie, err := r.Cookie(server.CookieName)
	if err != nil {
		return UserAuthorizedData(false)
	}

	var sessionId int
	sessionId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		println(err.Error())
		return UserAuthorizedData(false) // Добавить ошибку в Errors. Ну и вообще сделать Errors
	}

	session, ok := server.Sessions[sessionId]
	if !ok {
		println("Session#" + strconv.Itoa(sessionId) + " not found")
		return UserAuthorizedData(false)
	}

	userId := session.UserId
	user, err := db.GetUserById(server.DataBase, userId, server.UsersTableName)
	if err != nil {
		return JoinData(&map[string]interface{}{
			"Error": err.Error()},
			UserAuthorizedData(false),
		)
	}

	err = user.LoadReservationsFromDB(server.DataBase, server.ReservationsTableName)
	if err != nil {
		println(err.Error())
	}

	return JoinData(&map[string]interface{}{"AuthorizedUser": user}, UserAuthorizedData(true))
}

func (server *Server) GetTemplateAndUserData(dataToAdd []*map[string]interface{}, r *http.Request) *map[string]interface{} {
	return JoinData(
		JoinData(JoinDataArr(dataToAdd), server.BaseTemplateData),
		GetTemplateUserData(server, r))
}
