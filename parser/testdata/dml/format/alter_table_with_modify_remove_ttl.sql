-- Origin SQL:
ALTER TABLE infra.flow_processed_emails_local ON CLUSTER default_cluster MODIFY REMOVE TTL;

-- Format SQL:
ALTER TABLE infra.flow_processed_emails_local
ON CLUSTER default_cluster
MODIFY REMOVE TTL;
