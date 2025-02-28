-- Origin SQL:
CREATE MATERIALIZED VIEW infra_bm.view_name
    ON CLUSTER 'default_cluster' TO infra_bm.table_name
(
  `f1` DateTime64(3),
  `f2` String,
  `f3` String,
  `f4` String,
  `f5` String,
  `f6` Int64
) AS
SELECT f1,
       f2,
       visitParamExtractString(properties, 'f3') AS f3,
       visitParamExtractString(properties, 'f4') AS f4,
       visitParamExtractString(properties, 'f5') AS f5,
       visitParamExtractInt(properties, 'f6') AS f6
FROM infra_bm.table_name1
WHERE infra_bm.table_name1.event = 'test-event' AND
    NOT isZeroOrNull(f2) AND f6-2 > 0

-- Format SQL:
CREATE MATERIALIZED VIEW infra_bm.view_name ON CLUSTER 'default_cluster' TO infra_bm.table_name (`f1` DateTime64(3), `f2` String, `f3` String, `f4` String, `f5` String, `f6` Int64) AS SELECT f1, f2, visitParamExtractString(properties, 'f3') AS f3, visitParamExtractString(properties, 'f4') AS f4, visitParamExtractString(properties, 'f5') AS f5, visitParamExtractInt(properties, 'f6') AS f6 FROM infra_bm.table_name1 WHERE infra_bm.table_name1.event = 'test-event' AND NOT isZeroOrNull(f2) AND f6 - 2 > 0;
