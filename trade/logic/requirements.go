package logic

import (
	"SyncBoxi40/trade/database"
)

func AddRequirement(info, userid string) {
	println(info, userid)
	database.AddRequirement(info, userid)
}
func GetRequirements() (string, error) {
	return database.GetRequirements()
}

func DeleteRequirement(info string) {
	//println("Entro al logics")
	//println(strings.Split(info, ""))
	database.DeleteRequirement(info)
}

func Editrequirement(info string) {
	//println("Entro al logic editar")
	//println(info)
	database.Editrequirement(info)
}

// func RequirementsController(info string, action string) (string, error) {
// 	switch action {
// 	case "get":
// 		println("entro")
// 		return database.GetRequirements()
// 	case "delete":
// 		return database.GetCarteraCliente(info)
// 	default:
// 		return "", errors.New("invalid option")
// 	}
// }
