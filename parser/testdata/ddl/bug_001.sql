CREATE MATERIALIZED VIEW IF NOT EXISTS db.table
            ON CLUSTER 'default_cluster' TO db.table_mv
AS
SELECT
    event_ts,
    org_id,
    visitParamExtractString(properties, 'x') AS x,
    visitParamExtractString(properties, 'y') AS y,
    visitParamExtractString(properties, 'z') AS z,
    visitParamExtractString(properties, 'a') AS a,
    visitParamExtractString(properties, 'b') AS b,
    visitParamExtractString(properties, 'c') AS c,
    visitParamExtractString(properties, 'd') AS d,
    visitParamExtractInt(properties, 'e') AS e,
    visitParamExtractInt(properties, 'f') AS f
FROM db.table
WHERE db.table.event = 'hello';