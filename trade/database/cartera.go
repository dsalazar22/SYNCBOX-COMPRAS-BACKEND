package database

import (
	"fmt"

	"SyncBoxi40/libs"
)

func GetCarteraCliente(info string) (string, error) {
	consulta := fmt.Sprintf(`

		drop table if exists _result;
		drop table if exists _details_result;
		drop table if exists _final_result;
		
		select concat(oi.type_document,'-',oi."document") as "document", nit, name, coalesce(total*trm,0) total,phone,coalesce(quota,0) quota, coalesce((cast(now() as date) - deadline),0) days, 
			deadline, coalesce(c.limit_day,0) as limit_day into temp _result
		from trade.outstanding_invoices oi right join master.customers c on oi.customer_id = c.customers_id 
		where c.nit='%s' and c.deleted = false;
		
		select *, case when days<=30 then total else 0 end as "30_days", case when days>30 and days<=60 then total else 0 end as "31_to_60_days", 
			case when days>60 and days<=90 then total else 0 end as "61_to_90_days",
			case when days>90 then total else 0 end as "+90_days", case when limit_day<=days then 1 else 0 end portfolio_ok,
			(select sum((orders_details.requested_amount - orders_details.delivered_quantity) * coalesce(orders_details.sale_price,0))
			from trade.orders inner join trade.orders_details on orders.order_id = orders_details.orders_id
				inner join master.customers on orders.customer_id=customers.customers_id
				inner join sysconfig.status_modules on status_modules.status_id=orders_details.status_module_id
			where status_modules.is_active = true and customers.nit='%s') as orders
			into temp _details_result
		from _result;
		
		select nit, name,phone, max(days) term_days, quota, sum("30_days") "30_days", sum("31_to_60_days") "31_to_60_days", sum("61_to_90_days") "61_to_90_days",
			sum("+90_days") "mas_90_days", sum(portfolio_ok) portfolio_ok, sum(total) total, quota< (sum(total)+orders),orders,limit_day,
			
			(select array_to_json(array_agg(row_to_json(d)))
				from (
					select * from _details_result dr where dr.nit=d.nit
				)d) as details into temp _final_result
			
		from _details_result as d
		group by nit, name,phone, quota, orders,limit_day;

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select * from _final_result order by total desc 
		)d
	`, info, info)

	return libs.SendDB(consulta)
}

func GetCartera() (string, error) {
	consulta := `

		drop table if exists _result;
		drop table if exists _details_result;
		drop table if exists _final_result;
		
		select concat(oi.type_document,'-',oi."document") as "document", nit, name, total*trm total,phone,term_days, coalesce(quota,0) quota, (cast(now() as date) - deadline) days, 
			deadline, coalesce(c.limit_day,0) as limit_day into temp _result
		from trade.outstanding_invoices oi inner join master.customers c on oi.customer_id = c.customers_id ;
		
		select *, case when days<=30 then total else 0 end as "30_days", case when days>30 and days<=60 then total else 0 end as "31_to_60_days", 
			case when days>60 and days<=90 then total else 0 end as "61_to_90_days",
			case when days>90 then total else 0 end as "+90_days", case when limit_day<days then 1 else 0 end portfolio_ok into temp _details_result
		from _result;
		
		select nit, name,phone, term_days, quota, sum("30_days") "30_days", sum("31_to_60_days") "31_to_60_days", sum("61_to_90_days") "61_to_90_days",
			sum("+90_days") "mas_90_days", sum(portfolio_ok) portfolio_ok, sum(total) total, quota< sum(total),
			
			(select array_to_json(array_agg(row_to_json(d)))
				from (
					select * from _details_result dr where dr.nit=d.nit
				)d) as details into temp _final_result
			
		from _details_result as d
		group by nit, name,phone, term_days, quota;

		select array_to_json(array_agg(row_to_json(d)))
		from (
			select * from _final_result order by total desc 
		)d
	`

	return libs.SendDB(consulta)
}
