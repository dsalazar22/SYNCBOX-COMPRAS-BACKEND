package apps

import (
	"SyncBoxi40/produccionlinea/logic"
	"net/http"
)

func ConsultarCentrosTrabajo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept, Pragma, cache-control, expires")

	if r.Method == "GET" {
		result, err := logic.ConsultarCentrosTrabajo()

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(205)
			w.Write([]byte(err.Error()))
		}
	}
}
