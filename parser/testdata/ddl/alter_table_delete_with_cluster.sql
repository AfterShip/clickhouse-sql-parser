ALTER TABLE test.events ON CLUSTER 'default_cluster' DELETE WHERE id = 123 AND status = 'deleted';
