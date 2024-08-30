-- Origin SQL:
with t1 as (
    select 'value1' as value
    ), t2 as (
select 'value2' as value
    ), t3 as (
select 'value3' as value
    )
select
    t1.value as value1,
    t2.value as value2,
    t3.value as value3
from
    t1
        join t2 on true
        join t3
        join t4 on true
        join t5


-- Format SQL:
WITH t1 AS (SELECT 'value1' AS value), t2 AS (SELECT 'value2' AS value), t3 AS (SELECT 'value3' AS value) SELECT t1.value AS value1, t2.value AS value2, t3.value AS value3 FROM t1 JOIN t2 ON true JOIN t3 JOIN t4 ON true JOIN t5;
