package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func InsertDeliveries(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, '%s' creater, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(deliveries_header_id bigint, orders_details_id bigint, amount numeric, product_id bigint, ubications_id bigint);
		
		INSERT INTO trade.deliveries
		(deliveries_header_id,orders_details_id, amount, created, creater, logs)
		select deliveries_header_id,orders_details_id, amount, now(), creater, tlogs
		from tmpTable;
		
		insert into inventory.log_movement_product (product_id, ubication_id, amount, factor_converter,module,document_module_id, created, creater, typemove_id,is_income)
		select product_id, ubications_id, amount, 1,'deliveries', orders_details_id, now(), cast('%s' as uuid), (SELECT code FROM sysconfig.move_type where is_delivery limit 1),false
		from tmpTable;

		update trade.orders_details set delivered_quantity = (select sum(amount) from trade.deliveries as d where d.orders_details_id = tmpTable.orders_details_id)
		from tmpTable
		where tmpTable.orders_details_id = orders_details.orders_details_id;
	`, infoUser.Username, libs.GenerarLog(info, "delivery", infoUser), info, infoUser.IdUser)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func InsertDeliveriesHeader(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, '%s' creater, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id int8,transportation_company varchar(200));
		
		INSERT INTO trade.deliveries_header
		(orders_id,transportation_company,created,creater,state,logs)
		select orders_id,transportation_company,now(),creater,'AC', tlogs
		from tmpTable;
		
	`, infoUser.Username, libs.GenerarLog(info, "delivery", infoUser), info)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func UpdateDeliveriesHeader(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, '%s' creater, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id int8,transportation_company varchar(200));
		
		UPDATE trade.deliveries_header SET transportation_company=t.transportation_company
		from tmpTable t
		where deliveries_header.deliveries_header_id=t.deliveries_header_id;
		
	`, infoUser.Username, libs.GenerarLog(info, "delivery", infoUser), info)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func InvoiceDeliveriesHeader(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, '%s' creater, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(deliveries_header_id int8,transportation_company varchar(200), invoice varchar(20), 
			placa varchar(20), guia varchar(20));
		
		UPDATE trade.deliveries_header SET transportation_company=t.transportation_company, 
			invoice_number=invoice, state='FA', placa=t.placa, guia=t.guia
		from tmpTable t
		where deliveries_header.deliveries_header_id=t.deliveries_header_id;
		
	`, infoUser.Username, libs.GenerarLog(info, "delivery", infoUser), info)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func GetDeliveries(info string) (string, error) {
	consulta := fmt.Sprintf(`
		select array_to_json(array_agg(row_to_json(d)))
		from (
			select d.deliveries_header_id , concat('R',d.deliveries_header_id) code_header,  deliveries_id, orders_details_id, amount, 
				to_char(d.created, 'YYYY-MM-DD HH24:MI:SS') created, d.creater, placa, guia
			from trade.deliveries as d inner join trade.deliveries_header on d.deliveries_header_id=deliveries_header.deliveries_header_id
			where d.deliveries_header_id=%s
			order by d.created
		)d
	`, info)

	// fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func GetPendingInvoice() (string, error) {
	consulta := `
			drop table if exists _result_deliv;
			drop table if exists _finalresult;
			select orders_id, dh.deliveries_header_id, d.orders_details_id, sum(amount) amount, transportation_company, invoice_number, 
				dh.created, dh.creater,placa,guia into temp _result_deliv
			from trade.deliveries d right join trade.deliveries_header dh on dh.deliveries_header_id = d.deliveries_header_id
			where dh.state = 'AC'
			group by dh.deliveries_header_id, d.orders_details_id, transportation_company, dh.created, dh.creater, invoice_number,placa,guia;
			
			select rd.deliveries_header_id,concat(document_type,orders.order_code) order_code, customers.nit, customers.name,rd.transportation_company, rd.invoice_number, 
				products.code product_code, rd.placa, rd.guia,
				products.description product_description, coalesce(orders_details.sale_price,0) sale_price, 
				to_char(orders_details.requested_date, 'YYYY-MM-DD') as requested_date, 
				orders_details.deadline, orders_details.requested_amount,orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
				coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
				products.product_id, status_modules.code as code_status, families.buffer_days,
				customer_products.code as customer_code, customer_products.description as customer_description, is_national, products.standard_packing,
				rd2.amount delivery_amount, to_char(rd.created,'YYYY-MM-DD') as deliv_created, orders.shipping_address into temp _finalresult
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join master.products on products.product_id=orders_details.product_id
				inner join master.families on families.families_id=products.families_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
				left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
				inner join _result_deliv rd on rd.orders_id=orders_details.orders_id 
				inner join _result_deliv rd2 on rd2.orders_details_id=orders_details.orders_details_id;
			
				select array_to_json(array_agg(row_to_json(d)))
				from (
					select  deliveries_header_id, concat('R',deliveries_header_id) code_deliveries_header_id,order_code,nit,name,transportation_company, invoice_number,shipping_address, placa, guia,
					(select array_to_json(array_agg(row_to_json(d)))
						from (
							select distinct product_code,product_description,sale_price,delivery_amount --,total_price,box_count
							from _finalresult as fr
							where fr.deliveries_header_id = f.deliveries_header_id
					)d)
					from _finalresult as f
					group by deliveries_header_id,order_code,nit,name,transportation_company, invoice_number,shipping_address, placa, guia
				)d
	`

	// fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func GetInvoiceGenerate() (string, error) {
	consulta := `
			drop table if exists _result_deliv;
			drop table if exists _finalresult;
			select orders_id, dh.deliveries_header_id, d.orders_details_id, sum(amount) amount, transportation_company, invoice_number, placa, guia,
			dh.created, dh.creater into temp _result_deliv
			from trade.deliveries d right join trade.deliveries_header dh on dh.deliveries_header_id = d.deliveries_header_id
			where dh.state = 'FA'
			group by dh.deliveries_header_id, d.orders_details_id, transportation_company, dh.created, dh.creater, invoice_number, placa, guia;
			
			select rd2.deliveries_header_id,concat(document_type,orders.order_code) order_code, customers.nit, customers.name,rd2.transportation_company, rd2.invoice_number, products.code product_code, 
				products.description product_description, coalesce(orders_details.sale_price,0) sale_price, shipping_address,
				to_char(orders_details.requested_date, 'YYYY-MM-DD') as requested_date, 
				orders_details.deadline, orders_details.requested_amount,orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
				coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
				products.product_id, status_modules.code as code_status, families.buffer_days,
				customer_products.code as customer_code, customer_products.description as customer_description, is_national, products.standard_packing,
				rd2.amount delivery_amount, to_char(rd2.created,'YYYY-MM-DD') as deliv_created, rd2.placa, rd2.guia into temp _finalresult
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join master.products on products.product_id=orders_details.product_id
				inner join master.families on families.families_id=products.families_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
				--inner join _result_deliv rd on rd.orders_id=orders_details.orders_id 
				inner join _result_deliv rd2 on rd2.orders_details_id=orders_details.orders_details_id
				left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id;
			
				select array_to_json(array_agg(row_to_json(d)))
				from (
					select  deliveries_header_id, concat('R',deliveries_header_id) code_deliveries_header_id,order_code,nit,name,transportation_company, invoice_number, 
						placa, guia,shipping_address,deliv_created,
					(select array_to_json(array_agg(row_to_json(d)))
						from (
							select distinct product_code,product_description,sale_price,delivery_amount --,total_price,box_count
							from _finalresult as fr
							where fr.deliveries_header_id = f.deliveries_header_id
					)d)
					from _finalresult as f
					group by deliveries_header_id,order_code,nit,name,transportation_company, invoice_number,shipping_address,placa, guia,deliv_created
				)d
	`

	// fmt.Println(consulta)

	return libs.SendDB(consulta)
}
