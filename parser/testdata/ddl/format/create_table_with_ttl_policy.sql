-- Origin SQL:
CREATE TABLE tab
(
    d DateTime,
    a Int
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(d)
ORDER BY d
TTL d + INTERVAL 1 MONTH DELETE,
    d + INTERVAL 1 WEEK TO VOLUME 'aaa',
    d + INTERVAL 2 WEEK TO DISK 'bbb';


CREATE TABLE table_with_where
(
    d DateTime,
    a Int
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(d)
ORDER BY d
TTL d + INTERVAL 1 MONTH DELETE WHERE toDayOfWeek(d) = 1;

CREATE TABLE table_for_recompression
(
    d DateTime,
    key UInt64,
    value String
) ENGINE MergeTree()
ORDER BY tuple()
PARTITION BY key
TTL d + INTERVAL 1 MONTH RECOMPRESS CODEC(ZSTD(17)), d + INTERVAL 1 YEAR RECOMPRESS CODEC(LZ4HC(10))
SETTINGS min_rows_for_wide_part = 0, min_bytes_for_wide_part = 0;


-- Format SQL:
CREATE TABLE tab (d DateTime, a Int) ENGINE = MergeTree PARTITION BY toYYYYMM(d) TTL d + INTERVAL 1 MONTH DELETE, d + INTERVAL 1 WEEK TO VOLUME 'aaa', d + INTERVAL 2 WEEK TO DISK 'bbb' ORDER BY d;
CREATE TABLE table_with_where (d DateTime, a Int) ENGINE = MergeTree PARTITION BY toYYYYMM(d) TTL d + INTERVAL 1 MONTH DELETE WHERE toDayOfWeek(d) = 1 ORDER BY d;
CREATE TABLE table_for_recompression (d DateTime, key UInt64, value String) ENGINE = MergeTree() PARTITION BY key TTL d + INTERVAL 1 MONTH RECOMPRESS CODEC(ZSTD(17)), d + INTERVAL 1 YEAR RECOMPRESS CODEC(LZ4HC(10)) SETTINGS min_rows_for_wide_part=0, min_bytes_for_wide_part=0 ORDER BY tuple();
