-- Origin SQL:
SELECT
    *,
    CASE
        WHEN col2 = 'value1' THEN 'when1'
        WHEN col3 = 'value2' THEN 'when2'
        ELSE 'else'
    END as check_result
FROM table_name
WHERE col1 = '123456789'


-- Beautify SQL:
SELECT
  *,
  CASE
    WHEN col2 = 'value1' THEN 'when1'
    WHEN col3 = 'value2' THEN 'when2'
    ELSE 'else'
  END AS check_result
FROM
  table_name
WHERE
  col1 = '123456789';
