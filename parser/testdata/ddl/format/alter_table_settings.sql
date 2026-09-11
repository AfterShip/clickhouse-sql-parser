-- Origin SQL:
ALTER TABLE test.events_local MODIFY COLUMN f0 Int64 SETTINGS alter_sync = 2;
ALTER TABLE test.events_local ADD COLUMN f1 String, DROP COLUMN f2 SETTINGS alter_sync = 2;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' MATERIALIZE INDEX IF EXISTS created_at_idx SETTINGS mutations_sync = 2;


-- Format SQL:
ALTER TABLE test.events_local MODIFY COLUMN f0 Int64 SETTINGS alter_sync=2;
ALTER TABLE test.events_local ADD COLUMN f1 String, DROP COLUMN f2 SETTINGS alter_sync=2;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' MATERIALIZE INDEX IF EXISTS created_at_idx SETTINGS mutations_sync=2;
