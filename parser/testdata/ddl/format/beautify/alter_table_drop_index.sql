-- Origin SQL:
ALTER TABLE test.event_local ON CLUSTER 'default_cluster' DROP INDEX f1;

-- Beautify SQL:
ALTER TABLE test.event_local
ON CLUSTER 'default_cluster'
DROP INDEX f1;
