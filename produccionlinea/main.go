package main

import (
	"SyncBoxi40/libs"
	"SyncBoxi40/produccionlinea/apps"
	"SyncBoxi40/produccionlinea/logic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func getConn() {
	content, err := ioutil.ReadFile("app.connect")
	// var configConnection dto.ConfigConnection

	connStr := libs.ConfigConnection{
		Addr:    "127.0.0.1:6379",
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

	// fmt.Println(connStr)

	fmt.Println(connStr)

	libs.StrConn = connStr
}

func main() {
	getConn()
	fmt.Println(logic.ConsultarCentrosTrabajo())
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/consultarcts", apps.ConsultarCentrosTrabajo)

	http.ListenAndServe(":3500", router)
}
