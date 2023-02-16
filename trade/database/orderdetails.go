package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func InsertOrdersDetails(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id bigint, requested_amount numeric, requested_date date, deadline date, product_id bigint, 
				sale_price numeric, discount numeric, tax numeric, status_module_id integer, um varchar(50));
		
		INSERT INTO trade.orders_details
		(orders_id, requested_amount, delivered_quantity, requested_date, deadline, status_module_id,product_id,logs,sale_price,discount, tax,um)
		select orders_id , requested_amount,0, requested_date, coalesce(deadline,requested_date), status_module_id, product_id, tlogs, sale_price, discount,tax,um
		from tmpTable;
	`, libs.GenerarLog(info, "new_order", infoUser), info)

	return libs.SendDB(consulta)
}

func UpdateOrdersDetails(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_details_id bigint, orders_id bigint, requested_amount numeric, requested_date date, sale_price numeric, discount numeric, tax numeric, um varchar(50));
		
		UPDATE trade.orders_details
		SET requested_amount = tmp.requested_amount, 
		requested_date = tmp.requested_date, 
		sale_price = tmp.sale_price, 
		discount = tmp.discount, 
		tax = tmp.tax, 
		um = tmp.um,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where orders_details.orders_details_id = tmp.orders_details_id;

		select orders_id from tmpTable;
	`, libs.GenerarLog(info, "update_order", infoUser), info)

	return libs.SendDB(consulta)
}

func DeleteOrdersDetails(info string) (string, error) {

	consulta := fmt.Sprintf(`
		delete from trade.orders_details where orders_details_id = %s
	`, info)
	return libs.SendDB(consulta)
}

func GetOrderDetails(info string) (string, error) {

	consulta := fmt.Sprintf(`

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select o.orders_details_id, o.orders_id, p.product_id, p.code, p.description, sale_price, requested_amount, coalesce(discount,0) discount,requested_date,um,
				delivered_quantity,
				coalesce(tax,0) tax,((sale_price*requested_amount)-(((sale_price*requested_amount)*coalesce(discount,0))/100))+(((sale_price*requested_amount)*coalesce(tax,0))/100) as total
			from trade.orders_details o inner join sysconfig.status_modules sm on o.status_module_id = sm.status_id 
				inner join master.products p on p.product_id = o.product_id 
			where o.orders_id =%s and ((delivered_quantity=0 and sm.is_active = false) =false)
		)d;

	`, info)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}
