-- Origin SQL:
select cast(1 as Float64) as value;
select cast(1, 'Float64') as value;
select (1 as Float64) as value;
select 1::Float64 as value;
select cast(a + 1 as String) as value;
select cast(a + 1, 'String') as value;
select cast(-1 as Int8) as value;


-- Format SQL:
SELECT CAST(1 AS Float64) AS value;
SELECT CAST(1, 'Float64') AS value;
SELECT (1 AS Float64) AS value;
SELECT 1::Float64 AS value;
SELECT CAST(a + 1 AS String) AS value;
SELECT CAST(a + 1, 'String') AS value;
SELECT CAST(-1 AS Int8) AS value;
