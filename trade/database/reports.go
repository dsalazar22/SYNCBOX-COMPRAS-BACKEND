package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func GetGlobalOrdersReport() (string, error) {
	consulta := `
		drop table if exists tmp_orders;
		drop table if exists _resultcolor;

		select round(cast(coalesce(1.0-(DATE_PART('day', orders_details.deadline - case when status_modules.is_active then now() else orders_details.deadline end)/ families.buffer_days),0)*100.0 as numeric),2) as consumed_buffer, 
			customers.nit, customers.name customer_name,  to_char(orders.created, 'YYYY-MM-DD') order_created, products.code product_code, 
			products.description product_description, cast(coalesce(orders_details.sale_price,0) as numeric) sale_price, 
			to_char(orders_details.requested_date, 'YYYY-MM-DD') as requested_date, 
			orders_details.deadline, cast(orders_details.requested_amount as numeric) as requested_amount, orders_details.delivered_quantity, orders_details.requested_amount - orders_details.delivered_quantity as pending_amount, 
			coalesce(orders.document_customer,'') document_customer, families.description as families_description, orders_details.orders_details_id, orders_details.orders_id, 
			products.product_id, status_modules.code as code_status, families.buffer_days, concat(document_type,orders.order_code) order_code,
			customer_products.code as customer_code, customer_products.description as customer_description, is_national into temp tmp_orders
		from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
			inner join master.customers on orders.customer_id=customers.customers_id
			inner join master.products on products.product_id=orders_details.product_id
			inner join master.families on families.families_id=products.families_id
			inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			left join config.customer_products on customer_products.product_id=products.product_id and customer_products.customer_id=customers.customers_id
		where status_modules.is_active = true and orders.approved and orders_details.requested_amount - orders_details.delivered_quantity>0;

		select case when consumed_buffer >100 then 'dark' when consumed_buffer<=100 and consumed_buffer>66 
			then 'danger' when consumed_buffer<=66 and consumed_buffer>33 then 'warning' when consumed_buffer<=33 and consumed_buffer>0 then 'success' else 'info' end as color,
			to_char(case when deadline < cast(now() as date) then cast(now() as date) else deadline end, 'YYYY') aaaa,
			to_char(case when deadline < cast(now() as date) then cast(now() as date) else deadline end, 'YYYYMM') aaaamm,
			to_char(case when deadline < cast(now() as date) then cast(now() as date) else deadline end, 'WW') ww, * into temp _resultcolor
		from tmp_orders;

		select array_to_json(array_agg(row_to_json(d)))
				from ( select

			(select array_to_json(array_agg(row_to_json(d)))
				from (
				
					select color, round(sum(requested_amount*sale_price),2) totalprice, count(requested_amount) cantidad,
					(select array_to_json(array_agg(row_to_json(d)))
								from (select r.customer_name,  round(sum(r.requested_amount*r.sale_price),2) totalprice, count(r.requested_amount) cantidad, 
								(select array_to_json(array_agg(row_to_json(d)))
										from (select r2.*,  round(r2.requested_amount*r2.sale_price,2)
								from _resultcolor as r2
								where r.customer_name=r2.customer_name
								order by (r2.requested_amount*r2.sale_price) desc)d) as ordersdetails
						from _resultcolor as r
						where r.color=_resultcolor.color
						group by r.customer_name
						order by sum(r.requested_amount*r.sale_price) desc)d) customerdetails
					from _resultcolor
					group by color
				
				)d) detallecolor,
	
		(select array_to_json(array_agg(row_to_json(d)))
			from (

			select  aaaamm,round(sum(requested_amount*sale_price),2) totalprice,
				(select array_to_json(array_agg(row_to_json(d)))
					from (
						select r.*
						from _resultcolor as r
						where r.aaaamm=_resultcolor.aaaamm)d) pedidosmes,
						
				
				(select array_to_json(array_agg(row_to_json(d)))
					from (
						select rx.aaaamm, rx.ww,round(sum(rx.requested_amount*rx.sale_price),2) totalprice,
							(select array_to_json(array_agg(row_to_json(d)))
								from (
									select r.*
									from _resultcolor as r
									where r.aaaamm=rx.aaaamm and r.ww=rx.ww)d)
						from _resultcolor rx
						where _resultcolor.aaaamm=rx.aaaamm
						group by rx.aaaamm,rx.ww
						order by rx.aaaamm,rx.ww
						
					)d) detallesemana
			from _resultcolor
			group by aaaamm
			order by aaaamm)d) detalletiempo ,
						
				
			(select array_to_json(array_agg(row_to_json(d)))
				from (
					select rx.aaaa, rx.ww,round(sum(rx.requested_amount*rx.sale_price),2) totalprice,
						(select array_to_json(array_agg(row_to_json(d)))
							from (
								select r.*
								from _resultcolor as r
								where r.aaaa=rx.aaaa and r.ww=rx.ww)d)
					from _resultcolor rx
					group by rx.aaaa,rx.ww
					order by rx.aaaa,rx.ww
					
				)d) detallesemana
		)d
	`

	fmt.Println(consulta)

	return libs.SendDB(consulta)
}
