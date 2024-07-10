-- Origin SQL:
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX my_index(f0) TYPE minmax GRANULARITY 1024;


-- Format SQL:
ALTER TABLE test.events_local
ON CLUSTER 'default_cluster'
ADD INDEX my_index(f0) TYPE minmax GRANULARITY 1024;
