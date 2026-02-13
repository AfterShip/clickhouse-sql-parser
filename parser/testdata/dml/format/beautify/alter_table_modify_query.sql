-- Origin SQL:
ALTER TABLE test.some_mv ON CLUSTER cluster MODIFY QUERY SELECT field1, field2 FROM test.some_table WHERE count >= 3;

-- Beautify SQL:
ALTER TABLE test.some_mv
ON CLUSTER cluster
MODIFY QUERY SELECT
  field1,
  field2
FROM
  test.some_table
WHERE
  count >= 3;
