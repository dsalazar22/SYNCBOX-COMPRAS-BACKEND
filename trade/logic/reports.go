package logic

import "SyncBoxi40/trade/database"

func GetGlobalOrdersReport() (string, error) {
	return database.GetGlobalOrdersReport()
}
