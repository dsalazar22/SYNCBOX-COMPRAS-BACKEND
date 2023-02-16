package logic

import "SyncBoxi40/produccionlinea/database"

func ConsultarCentrosTrabajo() (string, error) {
	return database.ObtenerCentrosTrabajo()
}
