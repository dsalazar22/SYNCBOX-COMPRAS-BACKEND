package database

import (
	"SyncBoxi40/libs"
	"fmt"
)

//func AddRequirement(info string, userid string)
func AddRequirement(info, userid string) (string, error) {
	//println("Entro")
	consulta := fmt.Sprintf(`


	drop table if exists tmpTable;

	CREATE TEMP TABLE tmpTable AS
	select *
	from json_to_record('%s') as
	x(product_id int8, product_code text ,product_description text,quantity float, delivery_date timestamp);


	insert into trade.requirements  (product_id, created, quantity, delivery_date, product_description, product_code, is_active, creater, deleted)
	select product_id, now(), quantity, delivery_date, product_description, product_code, true, '%s', false 
	from tmpTable
	
			
			`, info, userid)

	//fmt.Println(consulta)
	//fmt.Println("USUARIO:", userid)
	return libs.SendDB(consulta)

}

func GetRequirements() (string, error) {
	consulta := `
	select array_to_json(array_agg(row_to_json(d)))
	from (
		select *
		from trade.requirements 
	)d;
	`

	return libs.SendDB(consulta)
}

func DeleteRequirement(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable;

	CREATE TEMP TABLE tmpTable AS
	select requirement_id 
		from json_to_record('%s') as
		x(requirement_id int);
	
	update trade.requirements
	set deleted = true, is_active = false
	where requirement_id = (select * from tmpTable) ;

	`, info)
	//fmt.Println(consulta)
	return libs.SendDB(consulta)

}

func Editrequirement(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable;

	CREATE TEMP TABLE tmpTable AS
	select *
		from json_to_record('%s') as
		x(requirement_id int, delivery_date date, quantity float, is_active boolean);
	
	update trade.requirements
		set delivery_date  = (select t.delivery_date from tmpTable t) ,  quantity = (select t.quantity from tmpTable t), is_active = (select t.is_active from tmpTable t)
		where requirement_id = (select t.requirement_id from tmpTable t) ; 
	`, info)

	//fmt.Println(consulta)
	return libs.SendDB(consulta)
}
