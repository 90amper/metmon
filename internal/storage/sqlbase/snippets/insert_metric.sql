insert into metmon.metrics (name, type)
values ($1, $2)
on conflict ON CONSTRAINT uniq_name do nothing
returning metric_id