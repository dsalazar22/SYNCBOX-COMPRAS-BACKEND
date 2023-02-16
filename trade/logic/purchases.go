package logic

import (
	"errors"

	"SyncBoxi40/trade/database"
)

func PurchasesController(info string, action string) (string, error) {
	switch action {
	case "select-active-true":
		return database.GetPurchasesActive("true")
	case "select-active-false":
		return database.GetPurchasesActive("false")
	case "get-families":
		//println("entro")
		return database.GetFamilies(info)
	case "report-new-delivery":
		return database.SaveReportNewDelivery(info)
	case "close-purchase-order":
		return database.ClosePurchaseOrder(info)
	case "get-delivery-reports":
		return database.GetDeliveryReports(info)
	case "get-number-purchase-order":
		return database.GetNumberPurchaseOrder()
	case "change-delivery-date":
		return database.ChangeDeliveryDate(info)
	default:
		return "", errors.New("Option Invalid")
	}
}

//ORDEN DE COMPRA

func GetRequirementsForPurchaseOrders() (string, error) {
	//println("Entro al logics")
	return database.GetRequirementsForPurchaseOrders()
}

func GetPurchaseOrders() (string, error) {
	//println("Entro al logics")
	return database.GetPurchaseOrders()
}

func AddPurchaseOrder(info string) (string, error) {
	//println("Entro al logics")
	println(info)
	return database.AddPurchaseOrder(info)
}
