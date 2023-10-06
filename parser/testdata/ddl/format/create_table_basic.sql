-- Origin SQL:
-- It's a short link events table
/**
    * @name Short link events
    * @description It's a short link events table
 */
CREATE TABLE IF NOT EXISTS test.events_local (
    f0 String,
    f1 String CODEC(ZSTD(1)),
    f2 VARCHAR(255),
    f3 Datetime,
    f4 Datetime,
    f5 Map(String,String),
    f6 String,
    f7 Nested (
        f70 UInt32,
        f71 UInt32,
        f72 DateTime,
        f73 Int64,
        f74 Int64,
        f75 String
    ),
    f8 Datetime DEFAULT now()
) ENGINE = MergeTree
PRIMARY KEY (f0, f1, f2)
PARTITION BY toYYYYMMDD(f3)
TTL f3 + INTERVAL 6 MONTH
ORDER BY (f1,f2,f3)

-- Format SQL:
CREATE TABLE IF NOT EXISTS test.events_local
(
  f0 String,
  f1 String CODEC(ZSTD(1)),
  f2 VARCHAR(255),
  f3 Datetime,
  f4 Datetime,
  f5 Map(String,String),
  f6 String,
  f7 Nested(
    f70 UInt32,
    f71 UInt32,
    f72 DateTime,
    f73 Int64,
    f74 Int64,
    f75 String),
  f8 Datetime DEFAULT now()
)
ENGINE = MergeTree
PRIMARY KEY (f0, f1, f2)
PARTITION BY toYYYYMMDD(f3)
TTL f3 + INTERVAL 6 MONTH
ORDER BY (f1, f2, f3);
