CREATE TABLE t (
    j JSON(message String, a.b UInt64, max_dynamic_paths=0, SKIP x, SKIP REGEXP 're')
) ENGINE = MergeTree
ORDER BY tuple();


