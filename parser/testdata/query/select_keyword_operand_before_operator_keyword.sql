SELECT CASE WHEN limit > 0 THEN limit ELSE offset END FROM quotas WHERE limit IN (1, 2) AND from = 1
