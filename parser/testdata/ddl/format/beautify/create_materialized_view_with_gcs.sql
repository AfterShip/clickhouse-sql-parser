-- Origin SQL:
CREATE MATERIALIZED VIEW database_name.view_name
        REFRESH EVERY 5 MINUTE TO database_name.table_name AS
        SELECT * FROM gcs(gcs_creds,url='https://storage.googleapis.com/some-bucket/some-path/');

-- Beautify SQL:
CREATE MATERIALIZED VIEW database_name.view_name REFRESH EVERY 5 MINUTE TO database_name.table_name AS SELECT
  *
FROM
  gcs(gcs_creds, url='https://storage.googleapis.com/some-bucket/some-path/');
