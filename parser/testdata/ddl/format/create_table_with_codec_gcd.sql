-- Origin SQL:
CREATE TABLE test_local
(
 `expectedrevenue` Nullable(Decimal(38, 4)) CODEC(GCD, ZSTD(1)),
 `count` UInt64 CODEC(GCD)
)
ENGINE = MergeTree
ORDER BY `count`;


-- Format SQL:
CREATE TABLE test_local (`expectedrevenue` Nullable(Decimal(38, 4)) CODEC(GCD, ZSTD(1)), `count` UInt64 CODEC(GCD)) ENGINE = MergeTree ORDER BY `count`;
