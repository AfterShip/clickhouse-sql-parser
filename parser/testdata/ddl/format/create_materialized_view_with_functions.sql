-- Origin SQL:
CREATE MATERIALIZED VIEW IF NOT EXISTS test.t0_view ON CLUSTER 'default_cluster' TO test.t0_name AS
SELECT 
	f1,
	f2,
	f3,
	attrs.1 as f4,
	arrayJoin(attrs.2) as f5
FROM
	test.t ARRAY
	JOIN [
		('string', mapKeys(string_attributes)),
		('int', mapKeys(int_attributes)),
		('float', mapKeys(float_attributes)),
		('bool', mapKeys(bool_attributes)),
		('null', mapKeys(null_attributes)),
		('int', 
			CAST(
				mapFilter(
					x -> NOT isZeroOrNull(x.1),
					map(
						'int1', `int1`,
						'int2', `int2`,
					)
				),
			'Map(String, String)')
		)
	] AS attrs
GROUP BY
	f1,
	f4

-- Format SQL:
CREATE MATERIALIZED VIEW IF NOT EXISTS test.t0_view ON CLUSTER 'default_cluster' TO test.t0_name AS SELECT f1, f2, f3, attrs.1 AS f4, arrayJoin(attrs.2) AS f5 FROM test.t  ARRAY JOIN [('string', mapKeys(string_attributes)), ('int', mapKeys(int_attributes)), ('float', mapKeys(float_attributes)), ('bool', mapKeys(bool_attributes)), ('null', mapKeys(null_attributes)), ('int', CAST(mapFilter(x -> -isZeroOrNull(x.1), map('int1', `int1`, 'int2', `int2`)), 'Map(String, String)'))] AS attrs GROUP BY f1, f4;
