-- Origin SQL:
ALTER TABLE test.events ON CLUSTER 'default_cluster' freeze partition '2023-07-18';;

-- Format SQL:
ALTER TABLE test.events ON CLUSTER 'default_cluster' FREEZE PARTITION '2023-07-18';
