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
