package logic

import (
	"errors"
	"fmt"

	"SyncBoxi40/libs"
	"SyncBoxi40/trade/database"
)

func HeaderOrdersControllers(info string, infoUser libs.InfoUser, action string) (string, error) {
	s := ""
	var e error
	switch action {
	case "get":
		s, e = database.GetOrder(info)
	case "insert":
		s, e = database.InsertOrders(info, infoUser)
	case "update":
		s, e = database.UpdateOrdersHeader(info, infoUser)
	case "update-notes":
		s, e = database.UpdateOrderNotes(info)
	case "update-transporter":
		s, e = database.UpdateTransporteNotes(info)
	case "approved":
		s, e = database.ApprovedOrder(info, infoUser)
	case "released":
		s, e = database.ReleasedOrder(info, infoUser)
	case "cancel":
		s, e = database.CancelOrder(info, infoUser)
	// case "delete":
	default:
		e = errors.New("Invalid Option")
	}

	return s, e
}

func DeliveriesController(info string, infoUser libs.InfoUser, action string) (string, error) {
	s := ""
	var e error
	switch action {
	case "get":
		s, e = database.GetDeliveries(info)
		fmt.Println(s)
	case "insert":
		s, e = database.InsertDeliveries(info, infoUser)
	case "insert-header":
		s, e = database.InsertDeliveriesHeader(info, infoUser)
	case "get-pending-invoice":
		s, e = database.GetPendingInvoice()
	case "update-invoice":
		s, e = database.UpdateDeliveriesHeader(info, infoUser)
	case "fa-invoice":
		s, e = database.InvoiceDeliveriesHeader(info, infoUser)
	case "get-invoice":
		s, e = database.GetInvoiceGenerate()
	default:
		e = errors.New("Invalid Option")
	}

	return s, e
}

func OrdersDetailControllers(info string, infoUser libs.InfoUser, action string) (string, error) {
	s := ""
	var e error
	switch action {
	case "get":
		s, e = database.GetOrderDetails(info)
	case "insert":
		s, e = database.InsertOrdersDetails(info, infoUser)
	case "update":
		s, e = database.UpdateOrdersDetails(info, infoUser)
	case "delete":
		s, e = database.DeleteOrdersDetails(info)
	// case "delete":
	default:
		e = errors.New("Invalid Option")
	}

	return s, e
}
