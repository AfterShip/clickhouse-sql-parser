-- Origin SQL:
SELECT t1.Timestamp FROM my_table t1 INNER JOIN my_other_table t2 ON t1.a=t2.b

-- Format SQL:
SELECT t1.Timestamp FROM my_table AS t1 INNER JOIN my_other_table AS t2 ON t1.a = t2.b;
