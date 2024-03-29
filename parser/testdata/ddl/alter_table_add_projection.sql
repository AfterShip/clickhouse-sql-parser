ALTER TABLE visits_order ADD PROJECTION user_name_projection (SELECT * ORDER BY user_name);
