-- Origin SQL:
CREATE MATERIALIZED VIEW db1.config_memory_v0
REFRESH EVERY 1 SECOND
(
    `schedule_id` String,
    `sample_rate` Int8,
    `start_at` DateTime64(9),
    `end_at` DateTime64(9),
    `created_at` DateTime64(9),
    `created_by` String,
    `properties` Map(String, String),
    `cluster_percentage` Float64,
    `schedule_filters` String
)
ENGINE = Memory
SETTINGS min_rows_to_keep = 250000, max_rows_to_keep = 500000
DEFINER = default SQL SECURITY DEFINER
COMMENT 'test comment'
AS SELECT
    schedule_id,
    sample_rate,
    start_at,
    end_at,
    created_at,
    created_by,
    properties,
    cluster_percentage,
    schedule_filters
FROM
(
    SELECT
        *,
        row_number() OVER (PARTITION BY schedule_filters ORDER BY created_at DESC) AS rn
    FROM db1.config_v0
    WHERE schedule_filters != '{}'
)
WHERE rn = 1


-- Beautify SQL:
CREATE MATERIALIZED VIEW db1.config_memory_v0
REFRESH EVERY 1 SECOND
(
  `schedule_id` String,
  `sample_rate` Int8,
  `start_at` DateTime64(9),
  `end_at` DateTime64(9),
  `created_at` DateTime64(9),
  `created_by` String,
  `properties` Map(String, String),
  `cluster_percentage` Float64,
  `schedule_filters` String
)
ENGINE = Memory
SETTINGS
  min_rows_to_keep=250000,
  max_rows_to_keep=500000
DEFINER = default
SQL SECURITY DEFINER
AS
  SELECT
    schedule_id,
    sample_rate,
    start_at,
    end_at,
    created_at,
    created_by,
    properties,
    cluster_percentage,
    schedule_filters
  FROM
    (SELECT
      *,
      row_number() OVER (PARTITION BY schedule_filters ORDER BY
        created_at DESC) AS rn
    FROM
      db1.config_v0
    WHERE
      schedule_filters != '{}')
  WHERE
    rn = 1
COMMENT 'test comment';
