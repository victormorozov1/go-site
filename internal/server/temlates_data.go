package server

import "golang.org/x/exp/maps"

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

//func (server *Server) GetTemplateAndUserData
