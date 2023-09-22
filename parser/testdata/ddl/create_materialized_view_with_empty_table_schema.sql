CREATE MATERIALIZED VIEW test.t0 on cluster default_cluster
ENGINE = ReplicatedAggregatingMergeTree('/clickhouse/{layer}-{shard}/test/t0', '{replica}')
PARTITION BY toYYYYMM(f0)
ORDER BY (f0)
POPULATE AS
select f0,f1,f2,coalesce(f0,f1) as f333
from
    (select
         f0,f1,f2,
         ROW_NUMBER() over(partition by f0 order by coalesce(f1,f2)) as rn
     from test.t
     where f3 in ('foo', 'bar', 'test')
       and env ='test'
    ) as tmp
where rn = 1;