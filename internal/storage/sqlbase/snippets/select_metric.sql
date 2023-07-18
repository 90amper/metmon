select cm."time", cm."name", cm."type", cm.cvalue, cm.gvalue
from metmon.current_metrics cm
where cm.name = $1 and cm.type = $2