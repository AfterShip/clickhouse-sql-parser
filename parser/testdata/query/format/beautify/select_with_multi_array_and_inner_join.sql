-- Origin SQL:
SELECT
    JSONExtractString(t3.props, 'value') AS value
FROM t1
    ARRAY JOIN JSONExtractArrayRaw(t1.props, 'arr1') AS a1
    INNER JOIN t2 ON t2.id = JSONExtractString(a1, 'id')
    ARRAY JOIN JSONExtractArrayRaw(t2.props, 'arr2') AS a2
    INNER JOIN t3 ON t3.id = JSONExtractString(a2, 'id')
WHERE value != '';

-- Beautify SQL:
SELECT
  JSONExtractString(t3.props, 'value') AS value
FROM
  t1
  ARRAY JOIN
    JSONExtractArrayRaw(t1.props, 'arr1') AS a1
  INNER JOIN
    t2 ON t2.id = JSONExtractString(a1, 'id')
  ARRAY JOIN
    JSONExtractArrayRaw(t2.props, 'arr2') AS a2
  INNER JOIN
    t3 ON t3.id = JSONExtractString(a2, 'id')
WHERE
  value != '';
