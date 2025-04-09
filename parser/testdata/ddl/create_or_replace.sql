-- It's a short link events table
/**
    * @name Short link events
    * @description It's a short link events table
 */
CREATE OR REPLACE TABLE IF NOT EXISTS test.events_local (
    f0 String,
    f1 String CODEC(ZSTD(1)),
    f2 VARCHAR(255),
) ENGINE = MergeTree
PRIMARY KEY (f0, f1, f2)
PARTITION BY toYYYYMMDD(f1)
TTL f1 + INTERVAL 6 MONTH
ORDER BY (f1,f2)
COMMENT 'Comment for table';

CREATE OR REPLACE VIEW IF NOT EXISTS my_view(col1 String, col2 String)
AS
SELECT
    id,
    name
FROM
    my_table;