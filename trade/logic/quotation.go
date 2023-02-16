package logic

import (
	"errors"

	"SyncBoxi40/libs"
	"SyncBoxi40/trade/database"
)

func HeaderQuotationController(info string, infoUser libs.InfoUser, action string) (string, error) {
	s := ""
	var e error
	switch action {
	case "get":
		s, e = database.GetQuotation(info)
	case "getall":
		s, e = database.GetActivesQuot()
	case "insert":
		s, e = database.InsertQuotation(info, infoUser)
	case "insert-details":
		s, e = database.InsertQuotationDetails(info, infoUser)
	case "get-details":
		s, e = database.GetQuotationDetails(info)
	case "add-version":
		s, e = database.AddVersion(info)
	case "remove-version":
		s, e = database.RemoveVersion(info)
	case "delete-details":
		s, e = database.DeleteQuotationDetails(info)
	case "active-version":
		s, e = database.ActiveVersion(info)
	case "update-details":
		s, e = database.UpdateQuotationDetails(info, infoUser)
	case "update":
		s, e = database.UpdateOrdersQuotation(info, infoUser)
	case "update-notes":
		s, e = database.UpdateQuotationNotes(info)
	case "update-transporter":
		s, e = database.UpdateQuotationTransporter(info)
	case "approved":
		s, e = database.ApprovedOrder(info, infoUser)
	case "released":
		s, e = database.ReleasedOrder(info, infoUser)
	case "cancel":
		s, e = database.CancelQuot(info, infoUser)
	// case "delete":
	default:
		e = errors.New("invalid option")
	}

	return s, e
}
