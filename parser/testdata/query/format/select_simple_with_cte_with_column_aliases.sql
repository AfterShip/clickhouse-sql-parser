-- Origin SQL:
WITH
    test(f1, f2, f3) AS (SELECT f4, f5, f6 FROM sales)
SELECT
    f1 AS new_f1,
    f2 AS new_f2,
    f3 AS new_f3
FROM
    test;


-- Format SQL:
WITH test(f1, f2, f3) AS (SELECT f4, f5, f6 FROM sales) SELECT f1 AS new_f1, f2 AS new_f2, f3 AS new_f3 FROM test;
