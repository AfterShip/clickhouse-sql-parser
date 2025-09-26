-- Origin SQL:
CREATE TABLE t (
    j JSON(message String, a.b UInt64, max_dynamic_paths=0, SKIP x, SKIP REGEXP 're')
) ENGINE = MergeTree
ORDER BY tuple();



-- Format SQL:
CREATE TABLE t (j JSON(max_dynamic_paths=0, message String, a.b UInt64, SKIP x,  SKIP REGEXP 're')) ENGINE = MergeTree ORDER BY tuple();


