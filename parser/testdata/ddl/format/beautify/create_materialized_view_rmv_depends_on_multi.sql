-- Origin SQL:
CREATE MATERIALIZED VIEW db1.attrs_rollup_v0
REFRESH EVERY 5 MINUTE
DEPENDS ON db1.metadata_mv_v0, db1.metadata_mv_v1
ENGINE = Memory
AS SELECT DISTINCT
    service_name,
    attr_key,
    'profiles' AS dataset,
    'string' AS attr_type
FROM db1.metadata_v0


-- Beautify SQL:
CREATE MATERIALIZED VIEW db1.attrs_rollup_v0
REFRESH EVERY 5 MINUTE
DEPENDS ON db1.metadata_mv_v0, db1.metadata_mv_v1
ENGINE = Memory
AS
  SELECT DISTINCT
    service_name,
    attr_key,
    'profiles' AS dataset,
    'string' AS attr_type
  FROM
    db1.metadata_v0;
