package logic

import (
	"errors"

	"SyncBoxi40/trade/database"
)

func ControllerCartera(info string, action string) (string, error) {
	switch action {
	case "get":
		return database.GetCartera()
	case "get-customer":
		return database.GetCarteraCliente(info)
	default:
		return "", errors.New("Error")
	}
}
