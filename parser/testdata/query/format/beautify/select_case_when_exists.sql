-- Origin SQL:
SELECT
    *,
    CASE
        WHEN EXISTS(SELECT 1
    FROM table_name
    WHERE col1 = '999999999')
        THEN 'then'
        ELSE 'else'
    END as check_result
FROM table_name
WHERE col1 = '123456789'


-- Beautify SQL:
SELECT
  *,
  CASE
    WHEN EXISTS(SELECT
      1
    FROM
      table_name
    WHERE
      col1 = '999999999') THEN 'then'
    ELSE 'else'
  END AS check_result
FROM
  table_name
WHERE
  col1 = '123456789';
