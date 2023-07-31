insert into metmon.gauges  (metric_id, gvalue)
values ((select metric_id from metmon.metrics where name = $1), $2);