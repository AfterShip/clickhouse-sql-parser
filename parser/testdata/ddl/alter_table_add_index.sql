ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX my_index(f0) TYPE minmax GRANULARITY 1024;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX api_id_idx api_id TYPE set(100) GRANULARITY 2;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX arr_idx arr TYPE bloom_filter(0.01) GRANULARITY 3;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX content_idx content TYPE tokenbf_v1(30720, 2, 0) GRANULARITY 1;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD INDEX output_idx output TYPE ngrambf_v1(3, 10000, 2, 1) GRANULARITY 2;
