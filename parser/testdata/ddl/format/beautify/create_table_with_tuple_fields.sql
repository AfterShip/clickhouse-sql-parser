-- Origin SQL:
CREATE TABLE t0 on cluster default_cluster
(
    `tup0` Tuple(),
    `tup1` Tuple(String, Int64),
    `tup2` Tuple(String, Tuple(String, String)),
    `tup3` Tuple(a String, cd Tuple(c String, d String))
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}', '{replica}')
ORDER BY (tup1, tup2, tup3)
SETTINGS index_granularity = 8192;


-- Beautify SQL:
CREATE TABLE t0 ON CLUSTER default_cluster
(
  `tup0` Tuple(),
  `tup1` Tuple(String, Int64),
  `tup2` Tuple(String, Tuple(String, String)),
  `tup3` Tuple(a String, cd Tuple(c String, d String))
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}', '{replica}')
ORDER BY
  (tup1, tup2, tup3)
SETTINGS
  index_granularity=8192;
