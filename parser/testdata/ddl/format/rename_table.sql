-- Origin SQL:
RENAME TABLE t1 TO t11;
RENAME TABLE t1 TO t11 ON CLUSTER 'default_cluster';
RENAME TABLE t1 TO t11, t2 TO t22;
RENAME TABLE t1 TO t11, t2 TO t22 ON CLUSTER 'default_cluster';


-- Format SQL:
RENAME TABLE t1 TO t11;
RENAME TABLE t1 TO t11
ON CLUSTER 'default_cluster';
RENAME TABLE t1 TO t11, t2 TO t22;
RENAME TABLE t1 TO t11, t2 TO t22
ON CLUSTER 'default_cluster';
