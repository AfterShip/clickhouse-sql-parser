-- Origin SQL:
ALTER TABLE db.t0 ON CLUSTER default_cluster
MODIFY TTL
    toDateTime(timestamp / 1000000000) + INTERVAL 30 DAY TO DISK 'gcs',
    toDateTime(timestamp / 1000000000) + INTERVAL 60 DAY;


-- Beautify SQL:
ALTER TABLE db.t0
ON CLUSTER default_cluster
MODIFY TTL toDateTime(timestamp / 1000000000) + INTERVAL 30 DAY TO DISK 'gcs', toDateTime(timestamp / 1000000000) + INTERVAL 60 DAY;
