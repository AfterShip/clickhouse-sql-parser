-- rename table
RENAME TABLE t1 TO t11;
RENAME TABLE t1 TO t11 ON CLUSTER 'default_cluster';
RENAME TABLE t1 TO t11, t2 TO t22;
RENAME TABLE t1 TO t11, t2 TO t22 ON CLUSTER 'default_cluster';
-- rename dictionary   
RENAME DICTIONARY t1 TO t11;
RENAME DICTIONARY t1 TO t11 ON CLUSTER 'default_cluster';
RENAME DICTIONARY t1 TO t11, t2 TO t22;
RENAME DICTIONARY t1 TO t11, t2 TO t22 ON CLUSTER 'default_cluster';
-- rename database
RENAME DATABASE t1 TO t11;
RENAME DATABASE t1 TO t11 ON CLUSTER 'default_cluster';
RENAME DATABASE t1 TO t11, t2 TO t22;
RENAME DATABASE t1 TO t11, t2 TO t22 ON CLUSTER 'default_cluster';
