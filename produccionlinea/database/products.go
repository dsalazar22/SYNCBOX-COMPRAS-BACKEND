package database

import (
	"SyncBoxi40/libs"
)

func ObtenerProductos() (string, error) {
	consulta := `
	select array_to_json(array_agg(row_to_json(d)))
	from (
		select code, description 
		from master.products p 
	)d	
	`
	return libs.SendDB(consulta)

}
