SELECT toFloat64(limit), max(limit), if(limit > 0, limit, 0) FROM quotas
