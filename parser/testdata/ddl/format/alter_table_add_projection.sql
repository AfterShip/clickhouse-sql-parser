-- Origin SQL:
ALTER TABLE visits_order
ADD PROJECTION  IF NOT EXISTS user_name_projection
(SELECT * GROUP BY user_name ORDER BY user_name) AFTER a.user_id;


-- Format SQL:
ALTER TABLE visits_order ADD PROJECTION IF NOT EXISTS user_name_projection (SELECT * GROUP BY user_name ORDER BY user_name) AFTER a.user_id;
