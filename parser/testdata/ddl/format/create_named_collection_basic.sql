-- Origin SQL:
CREATE NAMED COLLECTION IF NOT EXISTS servercore_s3_config
AS url = 'http://local-minio:9000/*',
access_key_id = 'minioadmin',
secret_access_key = 'minioadmin';


-- Format SQL:
CREATE NAMED COLLECTION IF NOT EXISTS servercore_s3_config AS url = 'http://local-minio:9000/*', access_key_id = 'minioadmin', secret_access_key = 'minioadmin';
