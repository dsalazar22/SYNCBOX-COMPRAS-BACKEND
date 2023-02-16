package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"SyncBoxi40/libs"
	"SyncBoxi40/trade/database"
)

type infoOrders struct {
	Orders    []OrdersInfo
	Inventory []InventoryProducts
}

type InventoryProducts struct {
	Code         string
	Amount_total float64
}

type OrdersInfo struct {
	Visible_Buffer       string  `json:"visible_buffer"`
	Consumed_buffer      float64 `json:"consumed_buffer"`
	Nit                  string  `json:"nit"`
	Customer_name        string  `json:"customer_name"`
	Order_created        string  `json:"order_created"`
	Product_code         string  `json:"product_code"`
	Product_description  string  `json:"product_description"`
	Sale_price           float64 `json:"sale_price"`
	Requested_date       string  `json:"requested_date"`
	Deadline             string  `json:"deadline"`
	Requested_amount     float64 `json:"requested_amount"`
	Delivered_quantity   float64 `json:"delivered_quantity"`
	Pending_amount       float64 `json:"pending_amount"`
	DocumentCustomer     string  `json:"document_customer"`
	DocumentType         string  `json:"document_type"`
	Families_description string  `json:"families_description"`
	Orders_details_id    int64   `json:"orders_details_id"`
	Orders_id            int64   `json:"orders_id"`
	Product_id           int64   `json:"product_id"`
	Code_status          string  `json:"code_status"`
	Buffer_days          int     `json:"buffer_days"`
	Available_inventory  float64 `json:"available_inventory"`
	Total_pending        float64 `json:"total_pending"`
	Order_code           string  `json:"order_code"`
	Customer_description string  `json:"customer_description"`
	Customer_code        string  `json:"customer_code"`
	IsNational           bool    `json:"is_national"`
	// CellVariant          string  `json:"_cellVariants"`
}

func GetOC() ([]byte, error) {
	result, err := database.GetOrdersActive("true")
	infoInventory := calcInventory(result)
	return infoInventory, err
}

func OrdersController(info string, infoUser libs.InfoUser, action string) ([]byte, error) {
	switch action {
	case "detail":
		result, err := database.GetOrderDetail(info)
		return []byte(result), err
	case "select-orders":
		result, err := database.GetProductOrders()
		infoInventory := calcInventory(result)
		return infoInventory, err
	case "select-active-true":
		result, err := database.GetOrdersActive("true")
		infoInventory := calcInventory(result)
		return infoInventory, err
	case "select-active-false":
		result, err := database.GetOrdersActive("false")
		infoInventory := calcInventory(result)
		return infoInventory, err
	case "order-edit":
		result, err := database.UpdateOrders(info)
		client := libs.RedisConnect(2)
		defer client.Close()
		client.Set("order-edit:"+result, info, 60*time.Second)
		return []byte(result), err
	case "released-orders":
		result, err := database.GetOrdersPendingReleased()
		return []byte(result), err
	case "approved-orders":
		result, err := database.GetOrdersPendingApproved()
		return []byte(result), err
	// case "delete":
	default:
		return nil, errors.New("Option Invalid")
	}
}

func UpdateOrderStatus(info string, infoUser libs.InfoUser) (string, error) {
	return database.UpdateOrderStatus(info)
}

func calcInventory(result string) []byte {
	var items []infoOrders
	inventory := make(map[string]float64)

	json.Unmarshal([]byte(result), &items)

	for _, v := range items[0].Inventory {
		inventory[v.Code] = v.Amount_total
	}

	for i, v := range items[0].Orders {
		items[0].Orders[i].Available_inventory = 0
		items[0].Orders[i].Visible_Buffer = fmt.Sprintf("%.2f", items[0].Orders[i].Consumed_buffer)

		// t, _ := time.Parse("2006-01-02", items[0].Orders[i].Order_created)
		items[0].Orders[i].Order_created = strings.Replace(items[0].Orders[i].Order_created, "T00:00:00", "", -1)
		if v.Pending_amount > 0 {
			if inventory[v.Product_code] > 0 {
				if (inventory[v.Product_code] - v.Pending_amount) > 0 {
					inventory[v.Product_code] = inventory[v.Product_code] - v.Pending_amount
					items[0].Orders[i].Available_inventory = v.Pending_amount
					// items[0].Orders[i].Pending_amount = 0

				} else {
					items[0].Orders[i].Available_inventory = inventory[v.Product_code]
					inventory[v.Product_code] = 0
				}
			}
			items[0].Orders[i].Total_pending = (items[0].Orders[i].Requested_amount - items[0].Orders[i].Delivered_quantity) - items[0].Orders[i].Available_inventory
		}
	}
	orders, _ := json.Marshal(items[0].Orders)

	return orders
}
