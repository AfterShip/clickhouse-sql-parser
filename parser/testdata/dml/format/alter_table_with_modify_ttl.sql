-- Origin SQL:
ALTER TABLE infra.flow_processed_emails_local ON CLUSTER default_cluster MODIFY TTL created_at + INTERVAL 3 YEAR;

-- Format SQL:
ALTER TABLE infra.flow_processed_emails_local ON CLUSTER default_cluster MODIFY TTL created_at + INTERVAL 3 YEAR;
