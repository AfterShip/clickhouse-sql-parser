-- Origin SQL:
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD COLUMN f1 String AFTER f0 SETTINGS alter_sync = 2;


-- Beautify SQL:
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD COLUMN f1 String AFTER f0
SETTINGS
  alter_sync=2;
