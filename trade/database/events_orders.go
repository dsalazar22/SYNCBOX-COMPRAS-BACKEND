package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func InsertOrders(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(order_code varchar(250), customer_id integer, document_customer varchar(250), document_type varchar(50), invoice_address varchar(2500), 
			shipping_address varchar(2500), released bool, multiple_deliveries bool, referral bool, currency varchar(50), trm numeric, approved bool, 
			is_national bool, consultant_id uuid,quotation_id bigint, order_notes text, transporter_notes text, invoice_city varchar, shipping_city varchar, 
			shipping_city_code varchar, invoice_city_code varchar);
		
		update trade.quotation set approved = true where quotation_id=(select quotation_id from tmpTable);

		INSERT INTO trade.orders
		(order_code, customer_id, document_customer, document_type, invoice_address, shipping_address, created, logs, multiple_deliveries, referral, currency, trm, approved, released,is_national, consultant_id,quotation_id, order_notes, transporter_notes,invoice_city,shipping_city,shipping_city_code,invoice_city_code)
		select order_code, customer_id, document_customer, document_type, invoice_address, shipping_address, now(), tlogs, multiple_deliveries, referral, currency, trm, approved, released,is_national,'%s',quotation_id, order_notes, transporter_notes,invoice_city,shipping_city, shipping_city_code,invoice_city_code
		from tmpTable;

		select lastval();
	`, libs.GenerarLog(info, "new_order", infoUser), info, infoUser.IdUser)

	return libs.SendDB(consulta)
}

func UpdateOrdersHeader(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id bigint, order_code varchar(250), document_customer varchar(250), document_type varchar(50), invoice_address varchar(2500), 
			shipping_address varchar(2500), released bool, multiple_deliveries bool, referral bool, currency varchar(50), trm numeric, approved bool, 
			consultant_id uuid, invoice_city varchar, shipping_city varchar,shipping_city_code varchar, invoice_city_code varchar);
		
		UPDATE trade.orders
		SET order_code = tmp.order_code, 
		document_customer = tmp.document_customer, 
		document_type = tmp.document_type, 
		invoice_address = tmp.invoice_address, 
		shipping_address = tmp.shipping_address, 
		released = tmp.released,
		invoice_city=tmp.invoice_city,
		shipping_city=tmp.shipping_city,
		multiple_deliveries = tmp.multiple_deliveries,
		referral = tmp.referral,
		currency = tmp.currency,
		trm = tmp.trm, 
		consultant_id = tmp.consultant_id,
		approved=tmp.approved,
		shipping_city_code = tmp.shipping_city_code, 
		invoice_city_code = tmp.invoice_city_code,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where orders.order_id = tmp.orders_id;

		select orders_id from tmpTable;
	`, libs.GenerarLog(info, "update_order", infoUser), info)

	return libs.SendDB(consulta)
}

func ApprovedOrder(info string, infoUser libs.InfoUser) (string, error) {

	consecut := `select value from sysconfig.syncbox s where code = 'consec' and "module" = 'orders'`
	c, _ := libs.SendDB(consecut)

	consulta := ""
	if c == "true" {
		consulta = fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id bigint, approved bool);
		
		update trade.orders
		set approved=true,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where orders.order_id = tmp.orders_id;
	`, libs.GenerarLog(info, "approved", infoUser), info)
	} else {
		consulta = fmt.Sprintf(`
			CREATE TEMP TABLE tmpTable AS
			select *, cast('%s' as jsonb) as tlogs
			from json_to_record('%s') as 
				x(orders_id bigint, approved bool);

			update trade.orders
			set approved=true,
			order_code = (select consecutive +1 from sysconfig.document_types dt where code = (select document_type from trade.orders o where o.order_id = (select orders_id from tmpTable))),
			logs = case when logs is null then '[]' else logs end || tlogs
			from tmpTable as tmp
			where orders.order_id = tmp.orders_id;
	
			update sysconfig.document_types set consecutive = consecutive +1 where code = (select document_type from trade.orders o where o.order_id = (select orders_id from tmpTable));
	
			select consecutive from sysconfig.document_types dt where code = (select document_type from trade.orders o where o.order_id = (select orders_id from tmpTable))
	
	`, libs.GenerarLog(info, "approved", infoUser), info)
	}

	return libs.SendDB(consulta)
}

func ReleasedOrder(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(orders_id bigint, released bool);
		
		update trade.orders
		set released=tmp.released,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where orders.order_id = tmp.orders_id;
	`, libs.GenerarLog(info, "released", infoUser), info)

	return libs.SendDB(consulta)
}

func CancelOrder(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(order_id bigint, status_id integer);
		
		update trade.orders_details
		set status_module_id=tmp.status_id,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where orders_details.orders_id = tmp.order_id;
	`, libs.GenerarLog(info, "canceled", infoUser), info)

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func GetOrder(info string) (string, error) {

	consulta := fmt.Sprintf(`

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select order_id, order_code, customers_id, document_customer, document_type, invoice_address, 
				shipping_address, released, customers.name, customers.nit, is_national, customers.discount,
				orders.consultant_id, approved,order_notes,transporter_notes, to_char(orders.created, 'YYYY-MM-DD') created,
				to_char((select max(od.deadline) from trade.orders_details od where od.orders_id = orders.order_id ),'YYYY-MM-DD') deadline,
				customers.phone, customers.principal_contact, username, customers.email,invoice_city,shipping_city,
				shipping_city_code,invoice_city_code
			from trade.orders inner join master.customers on orders.customer_id=customers.customers_id
				left join users.users u on u.iduser = orders.consultant_id 
			where orders.order_id=%s
		)d;

	`, info)
	return libs.SendDB(consulta)
}

func UpdateOrderNotes(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(onotes text, ordid int);
		
		update trade.orders set order_notes=onotes
		from tmpTable
		where order_id=ordid;

		select ordid from tmpTable
	`, status)

	return libs.SendDB(consulta)
}

func UpdateTransporteNotes(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(onotes text, ordid int);
		
		update trade.orders set transporter_notes=onotes
		from tmpTable
		where order_id=ordid;

		select ordid from tmpTable
	`, status)

	return libs.SendDB(consulta)
}
