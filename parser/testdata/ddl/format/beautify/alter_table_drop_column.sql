-- Origin SQL:
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' DROP COLUMN IF EXISTS f1;

-- Beautify SQL:
ALTER TABLE test.events_local
ON CLUSTER 'default_cluster'
DROP COLUMN IF EXISTS f1;
