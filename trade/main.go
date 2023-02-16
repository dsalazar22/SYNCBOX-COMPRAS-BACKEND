package main

import (
	"SyncBoxi40/libs"
	"SyncBoxi40/trade/logic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var trmValue float32

func commercialsorders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")

	vars := mux.Vars(r)
	event := vars["event"]

	ua := r.Header.Get("Authorization")
	infoUser := libs.GetUser(ua, "commercial_orders", event)

	if infoUser.IdUser != "" {
		b, _ := ioutil.ReadAll(r.Body)

		if len(b) > 0 {
			result, err := logic.OrdersController(string(b), infoUser, event)
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write(result)
			} else {
				w.WriteHeader(202)
				w.Write([]byte(err.Error()))
			}
		} else {
			w.WriteHeader(205)
			w.Write([]byte(""))
		}
	}
}

func updateorderstatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")

	ua := r.Header.Get("Authorization")
	infoUser := libs.GetUser(ua, "commercial_orders", "event")

	if infoUser.IdUser != "" {
		b, _ := ioutil.ReadAll(r.Body)

		if len(b) > 0 {
			result, err := logic.UpdateOrderStatus(string(b), infoUser)

			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(result))
			} else {
				w.WriteHeader(202)
				w.Write([]byte(err.Error()))
			}
		} else {
			w.WriteHeader(205)
			w.Write([]byte(""))
		}
	}
}

func commercialpurchases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	//println("entro al main")
	if r.Method == "POST" {

		vars := mux.Vars(r)
		event := vars["event"]

		//ua := r.Header.Get("Authorization")
		//infoUser := libs.GetUser(ua, "commercial_purchases", event)
		// fmt.Println(ua, infoUser)
		// if infoUser.IdUser != "" {
		//println("c")
		b, _ := ioutil.ReadAll(r.Body)

		// print(string(b))

		if len(b) > 0 {
			//println("b")
			result, err := logic.PurchasesController(string(b), event)

			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(result))
			} else {
				w.WriteHeader(202)
				w.Write([]byte(err.Error()))
			}
		} else {
			w.WriteHeader(205)
			w.Write([]byte(""))
		}
		//}
	}
}

// func commercialchases(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")

// 	if r.Method == "POST" {

// 		vars := mux.Vars(r)
// 		event := vars["event"]

// 		ua := r.Header.Get("Authorization")
// 		infoUser := libs.GetUser(ua, "commercial_chases", event)
// 		// fmt.Println(ua, infoUser)
// 		if infoUser.IdUser != "" {
// 			b, _ := ioutil.ReadAll(r.Body)

// 			// print(string(b))

// 			if len(b) > 0 {
// 				result, err := logic.chasesController(string(b), infoUser, event)

// 				if err == nil {
// 					w.WriteHeader(http.StatusOK)
// 					w.Write([]byte(result))
// 				} else {
// 					w.WriteHeader(202)
// 					w.Write([]byte(err.Error()))
// 				}
// 			} else {
// 				w.WriteHeader(205)
// 				w.Write([]byte(""))
// 			}
// 		}
// 	}
// }
func HeaderQuotationController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		vars := mux.Vars(r)
		event := vars["event"]
		b, _ := ioutil.ReadAll(r.Body)
		ua := r.Header.Get("Authorization")
		infoUser := libs.GetUser(ua, "commercial_chases", event)
		result, err := logic.HeaderQuotationController(string(b), infoUser, event)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func addHeaderCommercialOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		vars := mux.Vars(r)
		event := vars["event"]
		b, _ := ioutil.ReadAll(r.Body)
		ua := r.Header.Get("Authorization")
		infoUser := libs.GetUser(ua, "commercial_chases", event)
		result, err := logic.HeaderOrdersControllers(string(b), infoUser, event)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func addCommercialOrderDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		vars := mux.Vars(r)
		event := vars["event"]
		b, _ := ioutil.ReadAll(r.Body)
		ua := r.Header.Get("Authorization")
		infoUser := libs.GetUser(ua, "commercial_chases", event)
		result, err := logic.OrdersDetailControllers(string(b), infoUser, event)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func deliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		vars := mux.Vars(r)
		event := vars["event"]
		b, _ := ioutil.ReadAll(r.Body)
		ua := r.Header.Get("Authorization")
		infoUser := libs.GetUser(ua, "deliveries", event)
		result, err := logic.DeliveriesController(string(b), infoUser, event)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}
	}
}

//Requerimientos
// func requirementsController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
// 	if r.Method == "POST" {

// 		vars := mux.Vars(r)
// 		event := vars["event"]
// 		b, _ := ioutil.ReadAll(r.Body)
// 		//ua := r.Header.Get("Authorization")
// 		//infoUser := libs.GetUser(ua, "deliveries", "event")
// 		//result, err :=

// 		// println("Requerimientos")
// 		// println(b)
// 		// println(string(b))

// 		logic.RequirementsController(string(b), event)
// 		//logic.AddRequirement(string(b), infoUser.IdUser)

// 		// if err == nil {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 		// } else {
// 		// 	w.WriteHeader(202)
// 		// 	w.Write([]byte(err.Error()))
// 		// }
// 	}
// }

func getRequirements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		result, err := logic.GetRequirements()

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}
func addRequirements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")

	fmt.Println("ENTRO")

	if r.Method == "POST" {

		//fmt.Println("ENTRO AL CONDICIONAL")

		b, _ := ioutil.ReadAll(r.Body)
		//fmt.Println("b", string(b))
		ua := r.Header.Get("Authorization")
		fmt.Println("ua", ua)

		infoUser := libs.GetUser(ua, "deliveries", "event")
		//fmt.Println("infouser", infoUser)
		//fmt.Println(infoUser.IdUser)

		if infoUser.IdUser != "" {

			logic.AddRequirement(string(b), infoUser.IdUser)
		}

		// if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		// } else {
		// 	w.WriteHeader(202)
		// 	w.Write([]byte(err.Error()))
		// }
	}
}

func deleteRequirements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		b, _ := ioutil.ReadAll(r.Body)
		//ua := r.Header.Get("Authorization")
		//infoUser := libs.GetUser(ua, "deliveries", "event")
		//result, err :=

		//println("Borrar Requerimiento")
		//println(string(b))

		logic.DeleteRequirement(string(b))
		//logic.AddRequirement(string(b), infoUser.IdUser)

		// if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		// } else {
		// 	w.WriteHeader(202)
		// 	w.Write([]byte(err.Error()))
		// }
	}
}

func editRequirements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		b, _ := ioutil.ReadAll(r.Body)
		//ua := r.Header.Get("Authorization")
		//infoUser := libs.GetUser(ua, "deliveries", "event")
		//result, err :=

		//println("Editar Requerimiento")
		//println(string(b))

		logic.Editrequirement(string(b))
		//logic.AddRequirement(string(b), infoUser.IdUser)

		// if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		// } else {
		// 	w.WriteHeader(202)
		// 	w.Write([]byte(err.Error()))
		// }
	}
}

//ORDEN DE COMPRA
func getRequirementsForPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	//println("Entro al back")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		result, err := logic.GetRequirementsForPurchaseOrders()

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func addPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "POST" {

		b, _ := ioutil.ReadAll(r.Body)
		ua := r.Header.Get("Authorization")
		infoUser := libs.GetUser(ua, "deliveries", "event")
		//fmt.Println(infoUser.IdUser)
		if infoUser.IdUser != "" {
			//result, err :=

			//println("Orden de compra")
			//println(b)
			//println(string(b))

			logic.AddPurchaseOrder(string(b))
			//logic.AddRequirement(string(b), infoUser.IdUser)
		}

		// if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		// } else {
		// 	w.WriteHeader(202)
		// 	w.Write([]byte(err.Error()))
		// }
	}
}

func getPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	//println("Entro al back")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		result, err := logic.GetPurchaseOrders()

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func getValueTRM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		// vars := mux.Vars(r)
		// event := vars["m"]

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%f", trmValue)))
	}
}

func getCartera(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		vars := mux.Vars(r)
		event := vars["event"]
		nit := vars["nit"]

		result, err := logic.ControllerCartera(nit, event)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func getglobalordersreport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,  Authorization")
	if r.Method == "GET" {

		result, err := logic.GetGlobalOrdersReport()

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(202)
			w.Write([]byte(err.Error()))
		}

	}
}

func getConn() {
	content, err := ioutil.ReadFile("app.connect")
	// var configConnection dto.ConfigConnection

	connStr := libs.ConfigConnection{
		Addr:    "192.168.115.175:5432",
		Pass:    "",
		Db:      0,
		UserURL: "http://127.0.0.1:1705",
	}
	if err != nil {
		fileConfig, _ := json.Marshal(connStr)
		ioutil.WriteFile("app.connect", fileConfig, os.FileMode(777))
	} else {
		json.Unmarshal(content, &connStr)
	}
	fmt.Println(connStr)
	libs.StrConn = connStr

}

func getTRM() {
	type TRM struct {
		ID     string  `json:"__id"`
		Valor  float32 `json:"valor"`
		Unidad string  `json:"unidad"`
		Desde  string  `json:"vigenciadesde"`
		Hasta  string  `json:"vigenciahasta"`
	}

	type Content struct {
		Context string `json:"@odata.context"`
		InfoTRM []TRM  `json:"value"`
	}

	var content Content
	resp, err := http.Get("https://www.datos.gov.co/api/odata/v4/32sa-8pi3?$skiptoken=row-srez-nx28_3s47&$top=1000000")

	if err != nil {
		log.Println(err)
		trmValue = 0
	} else {

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &content)

		trmValue = content.InfoTRM[len(content.InfoTRM)-1].Valor

	}
}

func main() {
	//fmt.Println("HOLA")
	getConn()
	//fmt.Println(logic.GetRequirements())

	logic.GetOC()

	go func() {
		for {
			getTRM()
			time.Sleep(time.Duration(60) * time.Minute)
		}
	}()

	// var infoUser libs.InfoUser
	// logic.OrdersController("", infoUser, "select-active")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/commercialsorders/{event}", commercialsorders)
	router.HandleFunc("/updateorderstatus", updateorderstatus)
	router.HandleFunc("/comercialordercontroller/{event}", addHeaderCommercialOrder)
	router.HandleFunc("/comercialorderdetailcontroller/{event}", addCommercialOrderDetail)
	router.HandleFunc("/gettrm/{m}", getValueTRM)
	router.HandleFunc("/getoutstandinginvoices/{event}/{nit}", getCartera)
	router.HandleFunc("/deliveries/{event}", deliveries)
	router.HandleFunc("/globalordersreport", getglobalordersreport)

	//router.HandleFunc("/commercialchases/{event}", commercialchases)
	router.HandleFunc("/commercialpurchases/{event}", commercialpurchases)

	//REQUERIMIENTOS
	//router.HandleFunc("/requirementscontroller/{event}", requirementsController)
	router.HandleFunc("/requirements", getRequirements)
	router.HandleFunc("/addrequirements", addRequirements)
	router.HandleFunc("/deleterequirements", deleteRequirements)
	router.HandleFunc("/editrequirements", editRequirements)

	//ORDENES DE COMPRAS
	router.HandleFunc("/getrequirementsforpurchaseorders", getRequirementsForPurchaseOrders)
	router.HandleFunc("/addpurchaseorder", addPurchaseOrder)
	router.HandleFunc("/getpurchaseorders", getPurchaseOrders)

	router.HandleFunc("/quotcontroller/{event}", HeaderQuotationController)

	http.ListenAndServe(":3300", router)

}
