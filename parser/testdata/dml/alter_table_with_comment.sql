ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD COLUMN a.f1 String default '' comment 'test' ;
ALTER TABLE test.events_local ON CLUSTER 'default_cluster' ADD COLUMN hello String default '';
