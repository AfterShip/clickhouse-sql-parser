-- Origin SQL:
SELECT * FROM test_table SETTINGS additional_table_filters={'test_table': 'status = 1'};

SELECT * FROM test_table SETTINGS additional_table_filters={'test_table': 'value = \'test\''};

SELECT * FROM test_table SETTINGS additional_table_filters={'test_table': 'value = ''test'''};

SELECT * FROM test_table
SETTINGS additional_table_filters={'test_table': 'id IN (\'a\', \'b\') AND status = \'active\''}
FORMAT JSON;

SELECT number, x, y FROM (SELECT number FROM system.numbers LIMIT 5) f
ANY LEFT JOIN (SELECT x, y FROM table_1) s ON f.number = s.x
SETTINGS additional_table_filters={'system.numbers':'number != 3', 'table_1':'x != 2'};


-- Beautify SQL:
SELECT
  *
FROM
  test_table
SETTINGS
  additional_table_filters={'test_table': 'status = 1'};
SELECT
  *
FROM
  test_table
SETTINGS
  additional_table_filters={'test_table': 'value = \'test\''};
SELECT
  *
FROM
  test_table
SETTINGS
  additional_table_filters={'test_table': 'value = ''test'''};
SELECT
  *
FROM
  test_table
SETTINGS
  additional_table_filters={'test_table': 'id IN (\'a\', \'b\') AND status = \'active\''}
FORMAT JSON;
SELECT
  number,
  x,
  y
FROM
  (SELECT
    number
  FROM
    system.numbers
  LIMIT 5) AS f
  ANY LEFT JOIN
    (SELECT
      x,
      y
    FROM
      table_1) AS s ON f.number = s.x
SETTINGS
  additional_table_filters={'system.numbers': 'number != 3', 'table_1': 'x != 2'};
