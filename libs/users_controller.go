package libs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GenerarLog(infoDato string, action string, infoUser InfoUser) string {
	t := time.Now()
	currentDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	result := fmt.Sprintf(`[{"action":"%s", "idUser":"%s", "user":"%s", "email":"%s", "date":"%s", "content":%s}]`, action, infoUser.IdUser, infoUser.Username, infoUser.Email, currentDate, infoDato)
	return result //fmt.Sprintf(`[{"action":"%s", "idUser":"%s", "user":"%s", "email":"%s", "date":"%s", "content":"{%s}"}]`, action, infoUser.IdUser, infoUser.Username, infoUser.Email, currentDate, infoDato)
}

func GetUser(code string, modulo string, event string) InfoUser {

	var result InfoUser

	// requestURL := "https://account.syncbox.cloud/obtenerInfoUsuario"

	requestURL := StrConn.UserURL + "/obtenerInfoUsuario"

	var client http.Client

	req, err := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", code)

	if err != nil {
	}
	resp, err3 := client.Do(req)

	if err3 != nil {
		panic(err3)
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	// fmt.Println(body)

	json.Unmarshal(body, &result)

	return result
}
