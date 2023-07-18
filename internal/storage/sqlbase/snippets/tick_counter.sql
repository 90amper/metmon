update metmon.counters
set cvalue = cvalue + 1
where metric_id = (
	select metric_id
	from metmon.metrics
	where name = $1
)