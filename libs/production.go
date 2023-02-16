package libs

type ProductionOrders struct {
	OrderProductionProcessId int64    `json:"order_production_process_id"`
	ConsecutiveOrder         int      `json:"consecutive_order"`
	WSPlanningDefaultID      int      `json:"ws_planning_default_id"`
	ProductionOrderID        int      `json:"production_order_id"`
	CustomerOrderId          string   `json:"customer_order_id"`
	Name                     string   `json:"name"`
	Code                     string   `json:"code"`
	ProductsDescription      string   `json:"products_description"`
	Deadline                 string   `json:"deadline"`
	ProgrammedAmount         float64  `json:"programed_amount"`
	ProducedAmount           float64  `json:"produced_amount"`
	PendingAmount            float64  `json:"pending_amount"`
	StatusModulesCode        string   `json:"status_modules_code"`
	ActivityDescription      string   `json:"activity_description"`
	Sequence                 int      `json:"sequence"`
	WorkstationCode          string   `json:"workstation_code"`
	ProductId                int64    `json:"product_id"`
	BufferDays               float64  `json:"buffer_days"`
	ConsumeBuffer            float64  `json:"consume_buffer"`
	WorkstationGroups        string   `json:"workstation_groups"`
	ProductionPerHour        int64    `json:"production_per_hour"`
	MinutePrepatation        int64    `json:"minute_prepatation"`
	StartDate                string   `json:"start_date"`
	FinishDate               string   `json:"finish_date"`
	WorkstationId            int64    `json:"workstation_id"`
	LastWorkstationId        int64    `json:"last_workstation_id"`
	Workstations             []string `json:"workstations"`
	JobId                    int64    `json:"job_id"`
	ProgrammedJob            float64  `json:"programmed_job"`
	LastOperation            bool     `json:"last_operation"`
	GlobalProcess            bool     `json:"global_process"`
	Released                 bool     `json:"released"`
	WorkstationPlantID       int64    `json:"workstation_plant_id"`
	Priority                 int64    `json:"priority"`
	WorkstationGroup         string   `json:"workstation_group"`
	RequestedDate            string   `json:"requested_date"`
	PlanningDate             string   `json:"planning_date"`
	ProjectedDate            string   `json:"projected_date"`
	Messages                 int      `json:"messages"`
}

type OrdersProductionPrepared struct {
	TotalWorkstation []ProductionOrders `json:"total_workstation"`
	TotalOrders      []ProductionOrders `json:"total_orders"`
	TotalProcess     []ProductionOrders `json:"total_process"`
	TotalPlanned     []ProductionOrders `json:"total_planned"`
}
