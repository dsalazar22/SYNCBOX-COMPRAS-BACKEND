package main

import{

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