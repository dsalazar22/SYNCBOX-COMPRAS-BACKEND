package libs

import "time"

type CalendarResult struct {
	WorkstationId    int32     `json:"workstation_id"`
	DefaultStartHour string    `json:"default_hour"`
	StartDate        time.Time `json:"start_date"`
	FinishDate       time.Time `json:"finish_date"`
	TotalHours       int32     `json:"total_hours"`
	ProgramedHours   int32     `json:"programed_hours"`
	Type             string    `json:"type"`
	Comment          string    `json:"comment"`
	Weekday          string    `json:"day_of_week"`
	NumberWeekday    int       `json:"number_weekday"`
	NumberWeek       int       `json:"number_week"`
	Year             int       `json:"year"`
	Month            int       `json:"month"`
	NameMonth        string    `json:"name_month"`
}

type Products struct {
	ProductId   int64
	Code        string
	Description string
	Inventory   float64
}

type BillOfMaterials struct {
	ParentId        int64
	ChildId         int64
	AmountPerParent float64
}

var ProductsMap map[int64]Products

// type NodeProducts struct {
// 	Parent      *NodeProducts
// 	Product     []*NodeProducts
// 	ProductId   int64
// 	Code        string
// 	Description string
// 	Inventory   float64
// }
