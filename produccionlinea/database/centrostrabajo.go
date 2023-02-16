package database

import "SyncBoxi40/libs"

func ObtenerCentrosTrabajo() (string, error) {
	consulta := `
		select array_to_json(array_agg(row_to_json(d)))
		from(
			select workstation_id, code, description 
			from master.workstation w 
		)d
	`

	return libs.SendDB(consulta)
}
