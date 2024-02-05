WITH
    t1 AS
        (
            SELECT 1 AS value
    ),
    t2 AS
       (
SELECT 2 AS value
    )
SELECT *
FROM t1
         LEFT JOIN t2 ON true