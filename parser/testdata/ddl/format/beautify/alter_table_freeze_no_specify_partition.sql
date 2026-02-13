-- Origin SQL:
ALTER TABLE test.events ON CLUSTER 'default_cluster' freeze;

-- Beautify SQL:
ALTER TABLE test.events
ON CLUSTER 'default_cluster'
FREEZE;
