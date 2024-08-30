-- Origin SQL:
SELECT replica_name FROM system.ha_replicas UNION DISTINCT SELECT replica_name FROM system.ha_unique_replicas format JSON

-- Format SQL:
SELECT replica_name FROM system.ha_replicas UNION DISTINCT SELECT replica_name FROM system.ha_unique_replicas FORMAT JSON;
