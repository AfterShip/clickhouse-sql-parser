-- Origin SQL:
CREATE TABLE test_local
(
 `foo` Nullable(Decimal(38, 4)) CODEC(GCD, ZSTD(1)),
 `bar` UInt64 CODEC(GCD)
)
ENGINE = MergeTree
ORDER BY `bar`;


-- Beautify SQL:
CREATE TABLE test_local
(
  `foo` Nullable(Decimal(38, 4)) CODEC(GCD, ZSTD(1)),
  `bar` UInt64 CODEC(GCD)
)
ENGINE = MergeTree
ORDER BY
  `bar`;
