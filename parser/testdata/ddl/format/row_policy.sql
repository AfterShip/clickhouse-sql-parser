-- Origin SQL:
CREATE ROW POLICY example_policy ON db.example_table USING id = 42 TO example_role;
CREATE POLICY IF NOT EXISTS example_policy ON db.example_table, db.example_table_local USING id = toUInt64(getSetting('filter_id')) TO example_role;
CREATE ROW POLICY OR REPLACE p1, p2 ON db.example_table IN local_directory ON CLUSTER 'example-cluster' AS RESTRICTIVE USING has([1, 2], id) TO ALL EXCEPT admin, CURRENT_USER;
CREATE POLICY all_tables ON db.*, literal_table ON db.`*` USING 1 TO reader;
ALTER ROW POLICY IF EXISTS example_policy ON db.example_table USING id = 99 WITH CHECK id > 0 TO reader, `service-user`@'localhost';
ALTER POLICY example_policy ON db.example_table RENAME TO renamed AS PERMISSIVE USING NONE TO NONE;
DROP ROW POLICY IF EXISTS p1, p2 ON db.example_table FROM local_directory ON CLUSTER 'example-cluster';
SHOW CREATE ROW POLICY example_policy ON db.example_table, example_policy ON db.example_table_local;
SHOW CREATE ROW POLICIES ON db.*;
SHOW ROW POLICIES;


-- Format SQL:
CREATE ROW POLICY example_policy ON db.example_table FOR SELECT USING id = 42 TO example_role;
CREATE ROW POLICY IF NOT EXISTS example_policy ON db.example_table, example_policy ON db.example_table_local FOR SELECT USING id = toUInt64(getSetting('filter_id')) TO example_role;
CREATE ROW POLICY OR REPLACE p1 ON db.example_table, p2 ON db.example_table IN local_directory ON CLUSTER 'example-cluster' AS RESTRICTIVE FOR SELECT USING has([1, 2], id) TO ALL EXCEPT admin, CURRENT_USER;
CREATE ROW POLICY all_tables ON db.*, literal_table ON db.`*` FOR SELECT USING 1 TO reader;
ALTER ROW POLICY IF EXISTS example_policy ON db.example_table FOR SELECT USING id = 99 WITH CHECK id > 0 TO reader, `service-user`@'localhost';
ALTER ROW POLICY example_policy ON db.example_table RENAME TO renamed AS PERMISSIVE FOR SELECT USING NONE TO NONE;
DROP ROW POLICY IF EXISTS p1 ON db.example_table, p2 ON db.example_table FROM local_directory ON CLUSTER 'example-cluster';
SHOW CREATE ROW POLICY example_policy ON db.example_table, example_policy ON db.example_table_local;
SHOW CREATE ROW POLICIES ON db.*;
SHOW ROW POLICIES;
