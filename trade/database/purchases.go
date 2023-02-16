package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func GetPurchasesActive(status string) (string, error) {
	consulta := fmt.Sprintf(`
	
		DROP TABLE IF EXISTS tmp_purchases;

		select cast(coalesce(1.0-(DATE_PART('day', purchases_details.deadline - case when status_modules.is_active then now() else purchases_details.deadline end)/ families.buffer_days),0)*100.0 as numeric) as consumed_buffer, 
		supplier.nit, supplier.name supplier_name, purchases.created order_created, products.code product_code, products.description product_description, coalesce(purchases_details.price,0) sale_price, 
		purchases_details.requested_date, purchases_details.deadline, purchases_details.requested_amount,purchases_details.delivered_quantity, 
		purchases_details.requested_amount - purchases_details.delivered_quantity as pending_amount, 
		families.description as families_description, purchases_details.purchases_details_id, purchases_details.purchases_id, products.product_id, 
		status_modules.code as code_status, families.buffer_days, products.details into temp tmp_purchases
		from trade.purchases inner join trade.purchases_details on purchases.purchases_id = purchases_details.purchases_id
			inner join master.supplier on purchases.supplier_id=supplier.supplier_id
			inner join master.products on products.product_id=purchases_details.product_id
			inner join master.families on families.families_id=products.families_id
			inner join sysconfig.status_modules on status_modules.status_id=purchases_details.status_module_id
			left join config.supplier_products on supplier_products.product_id=products.product_id and supplier_products.supplier_id=supplier.supplier_id
		where status_modules.is_active = %s;

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select round(consumed_buffer,2) as visible_buffer,*
			from tmp_purchases
			order by consumed_buffer desc
		)d
	
	`, status)
	rs, err := libs.SendDB(consulta)

	// println(rs)

	return rs, err
}

//ORDEN DE COMPRA

func GetRequirementsForPurchaseOrders() (string, error) {
	consulta := (`

	select array_to_json(array_agg(row_to_json(d)))
	from (
		select r.*, p.price_list 
		from trade.requirements r inner join master.products p 
		on r.product_id = p.product_id 
	)d;
	`)
	//println(consulta)
	return libs.SendDB(consulta)

}

func GetPurchaseOrders() (string, error) {
	consulta := (`
	
	select array_to_json(array_agg(row_to_json(d)))
	from (
		select *
		from trade.purchase_orders
	)d;

	`)

	return libs.SendDB(consulta)
}

func AddPurchaseOrder(info string) (string, error) {
	consulta := fmt.Sprintf(`
       
	drop table if exists tmpTable;

	CREATE TEMP TABLE tmpTable AS
	select *
	from json_to_record('%s') as
	x(requirement_id int8, product_id int8, date_agreed timestamp, quantity float8, trading_value float8,
		supplier_nit text,quantity_delivered float8, total_sale float8,
		is_active boolean);


	insert into trade.purchases(supplier_id, created, purchase_code)
	select (select supplier_id from master.supplier where nit = (select supplier_nit from tmpTable)) , now(), (select cast(nextval('trade.purchases_purchases_id_seq') + 1 as varchar))
	from tmpTable;

	   insert into trade.purchases_details (purchases_id, product_id, requested_amount, delivered_quantity, requested_date,
		deadline, status_module_id ,price) 
	   select (select last_value from trade.purchases_purchases_id_seq), product_id, quantity, quantity_delivered, date_agreed, date_agreed, 12, total_sale
	   from tmpTable;

	update trade.requirements
	set is_active = false
	 where requirement_id = (select requirement_id from tmpTable);
	`, info)

	println(consulta)

	return libs.SendDB(consulta)
}

func GetFamilies(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable ;

	create temp table tmpTable as 
	select *
	from json_to_record ('%s') as
	x(product_id int8);

	select array_to_json(array_agg(row_to_json(d)))
	from(
	select f.*
	from master.products p inner join master.families f on p.families_id = f.families_id 
	where p.product_id = (select product_id  from tmpTable)
	) d;

	`, info)

	//println(consulta)

	return libs.SendDB(consulta)
}

func SaveReportNewDelivery(info string) (string, error) {
	consulta := fmt.Sprintf(`
			drop table if exists tmpTable;

			CREATE TEMP TABLE tmpTable AS
			select *
			from json_to_record('%s') as
			x(purchases_id int8 ,quantity_delivered float8, price float8);

			insert into trade.purchase_receipt_report  (purchases_id, delivered_amount, price, date_creation)
			select purchases_id, quantity_delivered, price, now()
			from tmpTable;

			update trade.purchases_details
			set delivered_quantity = (delivered_quantity + (select quantity_delivered from tmpTable))
			where purchases_id  = (select purchases_id  from tmpTable);
			`, info)

	println(consulta)

	return libs.SendDB(consulta)

}

func ClosePurchaseOrder(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable; 

	create temp table tmpTable as
	select *
	from json_to_record ('%s') as 
	x(purchase_id int8);

	update trade.purchases_details 
	set status_module_id = 10, close_date = now() 
	where purchases_id = (select purchase_id from tmpTable);

	`, info)

	//println(consulta)

	return libs.SendDB(consulta)
}

func GetDeliveryReports(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable; 

	create temp table tmpTable as
	select *
	from json_to_record ('%s') as 
	x(purchases_id int8);

	select array_to_json(array_agg(row_to_json(d)))
	from(
	select * from trade.purchase_receipt_report prr 
	where prr.purchases_id = (select * from tmpTable)
	) d;
	`, info)

	println(consulta)
	return libs.SendDB(consulta)
}

func GetNumberPurchaseOrder() (string, error) {
	consulta := (`
	select array_to_json(array_agg(row_to_json(d)))
	from(
	select (last_value (purchases_id) over () + 1 )as purchase_order_number
	from trade.purchases p      
	limit 1	
	) d;
	`)

	//println(consulta)

	return libs.SendDB(consulta)
}

func ChangeDeliveryDate(info string) (string, error) {
	consulta := fmt.Sprintf(`
	drop table if exists tmpTable; 

	create temp table tmpTable as
	select *
	from json_to_record ('%s') as 
	x(purchases_id int8, new_delivery_date date);

	update trade.purchases_details 
	set requested_date = (select new_delivery_date from tmpTable),
	deadline = (select new_delivery_date from tmpTable)
	where purchases_id = (select purchases_id from tmpTable)
	`, info)

	println(consulta)
	return libs.SendDB(consulta)
}
