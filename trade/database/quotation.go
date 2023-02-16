package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

//Requerimientos

func GetQuotation(info string) (string, error) {

	consulta := fmt.Sprintf(`

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select quotation_id, order_code, customers_id, document_customer, document_type, invoice_address, 
				shipping_address, invoice_city,shipping_city, released, customers.name, customers.nit, is_national, customers.discount,
				quotation.consultant_id, approved,order_notes,transporter_notes, to_char(quotation.created, 'YYYY-MM-DD') created,
				to_char((select max(od.deadline) from trade.quotation_details od inner join trade.quotation_version qv on qv.quotation_version_id = od.quotation_version_id where qv.quotation_id = quotation.quotation_id ),'YYYY-MM-DD') deadline,
				customers.phone, customers.principal_contact, username, customers.email,
				to_char(quotation.created + interval '1 days'*expiration_days, 'YYYY-MM-DD') duedate,
			(select array_to_json(array_agg(row_to_json(d)))
			from (
				select quotation_version_id, consec, state, deleted
				from trade.quotation_version qv 
				where qv.quotation_id = quotation.quotation_id and qv.deleted=false
				order by consec
			)d) as versions, quotation.deleted, expiration_days,shipping_city_code,invoice_city_code
			from trade.quotation inner join master.customers on quotation.customer_id=customers.customers_id
				left join users.users u on u.iduser = quotation.consultant_id 
			where quotation.quotation_id=%s
		)d;

	`, info)
	return libs.SendDB(consulta)
}

func InsertQuotation(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs, cast('%s' as uuid) consultant_id 
		from json_to_record('%s') as 
			x(order_code varchar(250), customer_id integer, document_customer varchar(250), document_type varchar(50), invoice_address varchar(2500), 
			shipping_address varchar(2500), released bool, multiple_deliveries bool, referral bool, currency varchar(50), trm numeric, approved bool, is_national bool, 
			quotation_id smallint, expiration_days int, invoice_city varchar, shipping_city varchar,shipping_city_code varchar,invoice_city_code varchar);
		
		INSERT INTO trade.quotation
		(order_code, customer_id, document_customer, document_type, invoice_address, shipping_address, created, logs, referral, currency, trm, approved, released,is_national, consultant_id, expiration_days,invoice_city,shipping_city,shipping_city_code,invoice_city_code)
		select order_code, customer_id, document_customer, document_type, invoice_address, shipping_address, now(), tlogs, referral, currency, trm, approved, released,is_national,consultant_id, expiration_days,invoice_city,shipping_city,shipping_city_code,invoice_city_code
		from tmpTable;

		update tmpTable set quotation_id = lastval();

		insert into trade.quotation_version (quotation_id, consec, state, deleted)
		select quotation_id,1,true,false
		from tmpTable;

		select quotation_id
		from tmpTable;
	`, libs.GenerarLog(info, "new_order", infoUser), infoUser.IdUser, info)

	return libs.SendDB(consulta)
}

func AddVersion(info string) (string, error) {
	consulta := fmt.Sprintf(`
		insert into trade.quotation_version (quotation_id, consec, state, deleted)
		values (%s,(select count(quotation_id) from trade.quotation_version where quotation_id=%s)+1,false,false);
	`, info, info)

	return libs.SendDB(consulta)
}

func ActiveVersion(info string) (string, error) {
	consulta := fmt.Sprintf(`
		update trade.quotation_version set state=false where quotation_id = (select quotation_id from  trade.quotation_version where quotation_version_id=%s);

		update trade.quotation_version set state=true where quotation_version_id=%s;

	`, info, info)

	return libs.SendDB(consulta)
}

func RemoveVersion(info string) (string, error) {
	consulta := fmt.Sprintf(`
		update trade.quotation_version set deleted=true where quotation_version_id=%s;
	`, info)

	return libs.SendDB(consulta)
}

func UpdateOrdersQuotation(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(quotation_id bigint, order_code varchar(250), document_customer varchar(250), document_type varchar(50), invoice_address varchar(2500), 
			shipping_address varchar(2500), released bool, multiple_deliveries bool, referral bool, currency varchar(50), trm numeric, approved bool, consultant_id uuid, 
			expiration_days int, invoice_city varchar, shipping_city varchar,shipping_city_code varchar,invoice_city_code varchar);
		
		UPDATE trade.quotation
		SET order_code = tmp.order_code, 
		document_customer = tmp.document_customer, 
		document_type = tmp.document_type, 
		invoice_address = tmp.invoice_address, 
		shipping_address = tmp.shipping_address, 
		invoice_city=tmp.invoice_city,
		shipping_city=tmp.shipping_city,
		released = tmp.released,
		referral = tmp.referral,
		currency = tmp.currency,
		trm = tmp.trm, 
		consultant_id = tmp.consultant_id,
		approved=tmp.approved,
		logs = case when logs is null then '[]' else logs end || tlogs,
		expiration_days = tmp.expiration_days,
		shipping_city_code = tmp.shipping_city_code,
		invoice_city_code = tmp.invoice_city_code
		from tmpTable as tmp
		where quotation.quotation_id = tmp.quotation_id;

		select quotation_id from tmpTable;
	`, libs.GenerarLog(info, "update_order", infoUser), info)

	return libs.SendDB(consulta)
}

func UpdateQuotation(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(status_id smallint, dline date, ordid int);
		
		update trade.quotation_details set status_module_id=status_id,deadline=dline
		from tmpTable
		where quotation_details_id=ordid;

		select ordid from tmpTable
	`, status)

	// println(consulta)

	return libs.SendDB(consulta)
}

func UpdateQuotationStatus(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(quotation_id int, quotation_version_id int2, document_customer varchar, multiple_deliveries bool);

		update production.production_order set deadline = dline
		from tmpTable
		where customer_order_id =ordid;

		select ordid from tmpTable

	`, status)

	println(consulta)

	return libs.SendDB(consulta)
}

func GetQuotationActive(status string) (string, error) {
	consulta := fmt.Sprintf(`
		
		select 
			customers.nit, customers.name customer_name,  to_char(quotation.created, 'YYYY-MM-DD') order_created, products.code product_code, 
			products.description product_description, coalesce(quotation_details.sale_price,0) sale_price, 
			to_char(quotation_details.requested_date, 'YYYY-MM-DD') as requested_date, 
			quotation_details.deadline, quotation_details.requested_amount,quotation_details.delivered_quantity, quotation_details.requested_amount - quotation_details.delivered_quantity as pending_amount, 
			coalesce(quotation.document_customer,'') document_customer, families.description as families_description, quotation_details.quotation_details_id, quotation_details.quotation_id, 
			products.product_id, families.buffer_days, concat(document_type,quotation.order_code) order_code,
			customer_products.code as customer_code, customer_products.description as customer_description, is_national into temp tmp_quotation,
			invoice_city,shipping_city,shipping_city_code,invoice_city_code
		from trade.quotation inner join trade.quotation_details on quotation.quotation_id = quotation_details.quotation_id
			inner join master.customers on quotation.customer_id=customers.customers_id
			inner join master.products on products.product_id=quotation_details.product_id
			inner join master.families on families.families_id=products.families_id
			left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
		where quotation.deleted = false and quotation.approved = false and quotation.canceled = false;

		select array_to_json(array_agg(row_to_json(x)))
		from (
			select
				(select array_to_json(array_agg(row_to_json(d)))
				from (
					select *
					from tmp_orders
				)d) as orders
		)x
	
	`, status)
	fmt.Println(consulta)
	return libs.SendDB(consulta)
}

func GetQuotationDetail(status string) (string, error) {
	consulta := fmt.Sprintf(`
		drop table if exists _result_deliv;
		select orders_id, dh.deliveries_header_id, d.orders_details_id, sum(amount) amount, transportation_company, dh.created, dh.creater into temp _result_deliv
		from trade.deliveries d right join trade.deliveries_header dh on dh.deliveries_header_id = d.deliveries_header_id
		where dh.state = 'AC'
		group by dh.deliveries_header_id, d.orders_details_id, transportation_company, dh.created, dh.creater;
		
		select array_to_json(array_agg(row_to_json(d)))
		from (
			select distinct cast(coalesce(1.0-(DATE_PART('day', orders_details.deadline - case when status_modules.is_active then now() else orders_details.deadline end)/ families.buffer_days),0)*100.0 as double precision) as consumed_buffer, 
				customers.nit, customers.name customer_name, orders.created order_created, products.code product_code, 
				products.description product_description, coalesce(orders_details.sale_price,0) sale_price, 
				to_char(orders_details.requested_date, 'YYYY-MM-DD') as requested_date, 
				orders_details.deadline, orders_details.requested_amount,orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
				coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
				products.product_id, status_modules.code as code_status, families.buffer_days, concat(document_type,orders.order_code) order_code,
				customer_products.code as customer_code, customer_products.description as customer_description, is_national, products.standard_packing,
				rd.deliveries_header_id,rd2.amount delivery_amount,rd.transportation_company, to_char(rd.created,'YYYY-MM-DD') as deliv_created, invoice_city,shipping_city
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join master.products on products.product_id=orders_details.product_id
				inner join master.families on families.families_id=products.families_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
				left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
				left join _result_deliv rd on rd.orders_id=orders_details.orders_id 
				left join _result_deliv rd2 on rd2.orders_details_id=orders_details.orders_details_id
			where lower(concat(document_type,orders.order_code)) = lower('%s') and orders.approved 
		
		)d
	`, status)

	return libs.SendDB(consulta)
}

func GetOrdersQuotationApproved() (string, error) {
	consulta := `
		select array_to_json(array_agg(row_to_json(x)))
		from (
			select orders.order_id, customers.nit, customers.name customer_name, orders.created order_created,
				coalesce(orders.document_customer,'') document_customer,
				orders.order_code,is_national
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join master.products on products.product_id=orders_details.product_id
				inner join master.families on families.families_id=products.families_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			where status_modules.is_active = true and coalesce(approved,false) = false  and coalesce(released,false) = true 
		)x
	`
	return libs.SendDB(consulta)
}

func InsertQuotationDetails(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(quotation_version_id bigint, requested_amount numeric, requested_date date, deadline date, product_id bigint, 
				sale_price numeric, discount numeric, tax numeric, status_module_id integer, um varchar(50));
		
		INSERT INTO trade.quotation_details
		(quotation_version_id, requested_amount, delivered_quantity, requested_date, deadline, status_module_id,product_id,logs,sale_price,discount, tax,um)
		select coalesce(quotation_version_id,1) , requested_amount,0, requested_date, coalesce(deadline,requested_date), status_module_id, product_id, tlogs, sale_price, discount,tax,um
		from tmpTable;
	`, libs.GenerarLog(info, "new_order", infoUser), info)

	return libs.SendDB(consulta)
}

func UpdateQuotationDetails(info string, infoUser libs.InfoUser) (string, error) {

	consulta := fmt.Sprintf(`
		CREATE TEMP TABLE tmpTable AS
		select *, cast('%s' as jsonb) as tlogs
		from json_to_record('%s') as 
			x(quotation_details_id bigint, orders_id bigint, requested_amount numeric, requested_date date, sale_price numeric, discount numeric, tax numeric, um varchar(50));
		
		UPDATE trade.quotation_details
		SET requested_amount = tmp.requested_amount, 
		requested_date = tmp.requested_date, 
		sale_price = tmp.sale_price, 
		discount = tmp.discount, 
		tax = tmp.tax, 
		um = tmp.um,
		logs = case when logs is null then '[]' else logs end || tlogs
		from tmpTable as tmp
		where quotation_details.quotation_details_id = tmp.quotation_details_id;

		select orders_id from tmpTable;
	`, libs.GenerarLog(info, "update_order", infoUser), info)

	return libs.SendDB(consulta)
}

func DeleteQuotationDetails(info string) (string, error) {

	consulta := fmt.Sprintf(`
		delete from trade.quotation_details where quotation_details_id = %s
	`, info)
	fmt.Println(consulta)
	return libs.SendDB(consulta)
}

func GetQuotationDetails(info string) (string, error) {

	consulta := fmt.Sprintf(`

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select o.quotation_details_id, o.quotation_version_id, p.product_id, p.code, p.description, sale_price, 
				requested_amount, coalesce(discount,0) discount,requested_date,um,
				delivered_quantity,deadline,status_module_id,
				coalesce(tax,0) tax,((sale_price*requested_amount)-(((sale_price*requested_amount)*coalesce(discount,0))/100))+(((sale_price*requested_amount)*coalesce(tax,0))/100) as total
			from trade.quotation_details o inner join master.products p on p.product_id = o.product_id 
			where o.quotation_version_id =%s 
		)d;

	`, info)

	// fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func GetActivesQuot() (string, error) {
	consulta := `
		drop table if exists _quot;

		select u.username, q.quotation_id,q.expiration_days , c.nit, c."name", order_code, (select "user"
				from jsonb_to_recordset(q.logs) as 
					x("user" text, "date" date)
			order by "date" asc limit 1), 
			(select "date"
				from jsonb_to_recordset(q.logs) as 
					x("user" text, "date" date)
			order by "date" asc limit 1), qv.state, qv.consec into temp _quot
		from trade.quotation q inner join trade.quotation_version qv on q.quotation_id =qv.quotation_id 
			inner join master.customers c on c.customers_id = q.customer_id
			inner join users.users u on u.iduser = q.consultant_id 
		where qv.deleted = false and q.deleted=false and q.approved=false;
		
		select array_to_json(array_agg(row_to_json(d)))
		from (
			select username, quotation_id, nit, "name", order_code, "user", "date", count(quotation_id) total_versions,expiration_days
			from _quot
			group by username, quotation_id, nit, "name", order_code, "user", "date",expiration_days
		)d
	`

	return libs.SendDB(consulta)
}

func CancelQuot(info string, infoUser libs.InfoUser) (string, error) {
	consulta := fmt.Sprintf(`
		UPDATE trade.quotation
		SET deleted = true,
			logs = case when logs is null then '[]' else logs end || '%s'
		where quotation_id = %s;
	`, libs.GenerarLog(info, "deleted_quot", infoUser), info)

	// fmt.Println(consulta)

	return libs.SendDB(consulta)
}

func UpdateQuotationNotes(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(onotes text, ordid int);
		
		update trade.quotation set order_notes=onotes
		from tmpTable
		where quotation_id=ordid;

		select ordid from tmpTable
	`, status)

	return libs.SendDB(consulta)
}

func UpdateQuotationTransporter(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(onotes text, ordid int);
		
		update trade.quotation set transporter_notes=onotes
		from tmpTable
		where quotation_id=ordid;

		select ordid from tmpTable
	`, status)

	return libs.SendDB(consulta)
}
