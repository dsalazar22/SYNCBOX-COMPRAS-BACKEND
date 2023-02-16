package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func UpdateOrders(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(status_id smallint, dline date, ordid int,odetailid bigint);
		
		update trade.orders_details set status_module_id=status_id,deadline=dline,close_date=now()
		from tmpTable
		where orders_details_id=odetailid;

		select ordid from tmpTable
	`, status)

	// println(consulta)

	return libs.SendDB(consulta)
}

func UpdateOrderStatus(status string) (string, error) {
	consulta := fmt.Sprintf(`

		CREATE TEMP TABLE tmpTable AS
		select * 
		from json_to_record('%s') as 
			x(status_id smallint, dline date, ordid int, odetailid bigint, pedidopadre bigint);

		update trade.orders_details set status_module_id=status_id,deadline=dline,close_date=now()
		from tmpTable
		where orders_details_id=odetailid;
		
		update tmpTable set pedidopadre = (select product_id from trade.orders_details where orders_details_id=odetailid);

		update production.production_order set deadline = dline
		from tmpTable
		where customer_order_id =ordid and product_id=pedidopadre;

		select ordid from tmpTable;

	`, status)

	// println(consulta)

	return libs.SendDB(consulta)
}

func GetOrdersActive(status string) (string, error) {
	consulta := fmt.Sprintf(`
	--DROP INDEX trade.orders_details_status_module_id_idx;
	--DROP INDEX trade.orders_details_orders_id_idx;
	--DROP INDEX trade.orders_details_product_id_idx;

	--CREATE INDEX orders_details_orders_id_idx ON trade.orders_details (orders_id);
	--CREATE INDEX orders_details_product_id_idx ON trade.orders_details (product_id);
	--CREATE INDEX orders_details_status_module_id_idx ON trade.orders_details (status_module_id);
	

		
		select cast(coalesce(1.0-(DATE_PART('day', orders_details.deadline - case when status_modules.is_active then now() else orders_details.deadline end)/ families.buffer_days),0)*100.0 as double precision) as consumed_buffer, 
			customers.nit, customers.name customer_name,  to_char(orders.created, 'YYYY-MM-DD') order_created, products.code product_code, 
			products.description product_description, coalesce(orders_details.sale_price,0) sale_price, 
			to_char(orders_details.requested_date, 'YYYY-MM-DD') as requested_date, 
			orders_details.deadline, orders_details.requested_amount,orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
			coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
			products.product_id, status_modules.code as code_status, families.buffer_days, coalesce(orders.order_code, cast(orders.order_id as varchar(50))) order_code,
			customer_products.code as customer_code, customer_products.description as customer_description, is_national, orders.approved into temp tmp_orders
		from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
			inner join master.customers on orders.customer_id=customers.customers_id
			inner join master.products on products.product_id=orders_details.product_id
			inner join master.families on families.families_id=products.families_id
			inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
		where status_modules.is_active = %s and coalesce(products.service,false)=false and orders.released;

		select array_to_json(array_agg(row_to_json(x)))
		from (
			select
				(select array_to_json(array_agg(row_to_json(d)))
				from (
					select *
					from tmp_orders
					order by consumed_buffer desc
				)d) as orders,
	
				(select array_to_json(array_agg(row_to_json(d)))
				from (
					select code,amount_total
					from master.products inner join inventory.product_totals on product_totals.product_id=products.product_id
				)d) as inventory
		)x
	`, status)
	fmt.Println(consulta)
	return libs.SendDB(consulta)
}

func GetOrderDetail(status string) (string, error) {
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
				rd.deliveries_header_id,rd2.amount delivery_amount,rd.transportation_company, to_char(rd.created,'YYYY-MM-DD') as deliv_created
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

func GetOrdersPendingReleased() (string, error) {
	consulta := `
		select array_to_json(array_agg(row_to_json(x)))
		from (
			select DISTINCT u.username,orders.order_id, customers.nit, customers.name customer_name, orders.created order_created,
				coalesce(orders.document_customer,'') document_customer,
				orders.order_code,is_national
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join master.products on products.product_id=orders_details.product_id
				inner join master.families on families.families_id=products.families_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
				inner join users.users u on u.iduser = orders.consultant_id
			where status_modules.is_active = true and coalesce(released,false) = false
		)x
	`
	return libs.SendDB(consulta)
}

func GetOrdersPendingApproved() (string, error) {
	consulta := `
		drop table if exists _result;
		drop table if exists _orders;
		
		select oi.customer_id, sum(coalesce(total*trm,0)) total, coalesce(quota,0) quota, max(coalesce((cast(now() as date) - deadline),0)) days, 
			coalesce(c.limit_day,0) as limit_day into temp _result
		from trade.outstanding_invoices oi right join master.customers c on oi.customer_id = c.customers_id 
		where c.deleted = false
		group by oi.customer_id, coalesce(quota,0),coalesce(c.limit_day,0)
		having  sum(coalesce(total*trm,0)) >0;


		select orders.customer_id, sum((orders_details.requested_amount - orders_details.delivered_quantity) * coalesce(orders_details.sale_price,0)) as total into temp _orders
		from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
			inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
		where status_modules.is_active = true 
		group by orders.customer_id;	

		select array_to_json(array_agg(row_to_json(x)))
		from (
			select distinct orders.order_id, customers.nit, customers.name customer_name, to_char(orders.created, 'YYYY-MM-DD HH24:MI') order_created,
			coalesce(orders.document_customer,'') document_customer, coalesce(days,0) days, customers.limit_day,
			orders.order_code,is_national, _orders.total as orders,  coalesce(_result.total,0) invoices, customers.quota
		from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
			inner join master.customers on orders.customer_id=customers.customers_id
			inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			left join _result on _result.customer_id=customers.customers_id
			left join _orders on _orders.customer_id=customers.customers_id
		where status_modules.is_active = true and coalesce(approved,false) = false  and coalesce(released,false) = true 
		)x
	`
	return libs.SendDB(consulta)
}

func GetProductOrders() (string, error) {
	consulta := //fmt.Sprintf(
		`
		--drop table tmp_orders;
		--drop table _production;

		select DISTINCT production_order.production_order_id, production_order.consecutive_order, orders.order_id, production_order.product_id into temp _production
		from production.production_order inner join trade.orders on production_order.customer_order_id=orders.order_id
			inner join trade.orders_details on orders_details.orders_id=orders.order_id and orders_details.product_id=production_order.product_id
			inner join sysconfig.status_modules on status_modules.status_id=production_order.document_status_id
		where status_modules.is_active;

		select cast(coalesce(1.0-(DATE_PART('day', orders_details.deadline - case when status_modules.is_active then now() else orders_details.deadline end)/ families.buffer_days),0)*100.0 as double precision) as consumed_buffer, 
			customers.nit, customers.name customer_name, orders.created order_created, products.code product_code, 
			products.description product_description, coalesce(orders_details.sale_price,0) sale_price, orders_details.requested_date, 
			orders_details.deadline, orders_details.requested_amount,orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
			coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
			products.product_id, status_modules.code as code_status, families.buffer_days, orders.order_code, is_national into temp tmp_orders
		from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
			inner join master.customers on orders.customer_id=customers.customers_id
			inner join master.products on products.product_id=orders_details.product_id
			inner join master.families on families.families_id=products.families_id
			inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
		where status_modules.is_active;
		

		select array_to_json(array_agg(row_to_json(x)))
		from (
			select
				(select array_to_json(array_agg(row_to_json(d)))
				from (
					select tmp_orders.* 
					from tmp_orders left join _production on tmp_orders.orders_id= _production.order_id and tmp_orders.product_id= _production.product_id
					where _production.product_id is null and pending_amount>0
					order by consumed_buffer desc
				)d) as orders,
	
				(select array_to_json(array_agg(row_to_json(d)))
				from (
					select code,amount_total
					from master.products inner join inventory.product_totals on product_totals.product_id=products.product_id
				)d) as inventory
		)x
	
	` //, info)

	return libs.SendDB(consulta)
}
