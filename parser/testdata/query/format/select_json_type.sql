-- Origin SQL:
SELECT a, a.b, a.b.c.d.e;
SELECT JSON_TYPE('{"a": 1, "b": {"c": 2}}', '$.b');
SELECT CAST(some, 'String') AS value;
SELECT CAST(some.long, 'String') AS value;
SELECT CAST(some.long.json, 'String') AS value;
SELECT CAST(some.long.json.path, 'String') AS value;



-- Format SQL:
SELECT a, a.b, a.b.c.d.e;
SELECT JSON_TYPE('{"a": 1, "b": {"c": 2}}', '$.b');
SELECT CAST(some, 'String') AS value;
SELECT CAST(some.long, 'String') AS value;
SELECT CAST(some.long.json, 'String') AS value;
SELECT CAST(some.long.json.path, 'String') AS value;
