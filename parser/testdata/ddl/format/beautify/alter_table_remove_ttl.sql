-- Origin SQL:
ALTER TABLE test.events ON CLUSTER 'default_cluster' REMOVE TTL;

-- Beautify SQL:
ALTER TABLE test.events
ON CLUSTER 'default_cluster'
REMOVE TTL;
