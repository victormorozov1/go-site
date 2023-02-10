package server

import "golang.org/x/exp/maps"

func JoinData(data1 *map[string]interface{}, data2 *map[string]interface{}) *map[string]interface{} { // Вынести куда-нибудь
	resData := make(map[string]interface{})
	maps.Copy(resData, *data1)
	maps.Copy(resData, *data2)
	return &resData
}

func AddRoutesData(data *map[string]interface{}, server *Server) {
	(*data)["Routes"] = server.Routes
}
