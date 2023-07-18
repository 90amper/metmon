insert into metmon.counters (metric_id, cvalue)
values ((select metric_id from metmon.metrics where name = $1), $2)
on conflict (metric_id) do
update set
cvalue = excluded.cvalue,
"timestamp" = now();