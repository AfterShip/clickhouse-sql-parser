-- Origin SQL:
SELECT a FROM t WHERE interval > 1;
SELECT a FROM t WHERE interval = 1;
SELECT a FROM t WHERE interval IS NULL;
SELECT a FROM t WHERE interval IN (1, 2);
SELECT count() FROM t WHERE interval BETWEEN 1 AND 2;
SELECT a FROM t PREWHERE interval > 1;
SELECT a FROM t HAVING interval > 1;
SELECT a FROM t GROUP BY (interval);
SELECT a FROM t GROUP BY (a, interval);
SELECT max(interval) FROM t;
SELECT sum(interval) / count() FROM t;
SELECT if(interval > 1, 1, 2) FROM t;
SELECT interval + 1 FROM t;
SELECT interval * 2 FROM t;
SELECT interval[1] FROM t;
SELECT 1 AS interval, interval + 1;
SELECT a FROM t ORDER BY interval ASC;
SELECT a FROM t ORDER BY interval DESC;
SELECT a FROM t ORDER BY b ASC, interval DESC;
SELECT a FROM t ORDER BY interval WITH FILL;
SELECT a FROM t ORDER BY interval WITH FILL STEP 1;


-- Format SQL:
SELECT a FROM t WHERE interval > 1;
SELECT a FROM t WHERE interval = 1;
SELECT a FROM t WHERE interval IS NULL;
SELECT a FROM t WHERE interval IN (1, 2);
SELECT count() FROM t WHERE interval BETWEEN 1 AND 2;
SELECT a FROM t PREWHERE interval > 1;
SELECT a FROM t HAVING interval > 1;
SELECT a FROM t GROUP BY (interval);
SELECT a FROM t GROUP BY (a, interval);
SELECT max(interval) FROM t;
SELECT sum(interval) / count() FROM t;
SELECT if(interval > 1, 1, 2) FROM t;
SELECT interval + 1 FROM t;
SELECT interval * 2 FROM t;
SELECT interval[1] FROM t;
SELECT 1 AS interval, interval + 1;
SELECT a FROM t ORDER BY interval ASC;
SELECT a FROM t ORDER BY interval DESC;
SELECT a FROM t ORDER BY b ASC, interval DESC;
SELECT a FROM t ORDER BY interval WITH FILL;
SELECT a FROM t ORDER BY interval WITH FILL STEP 1;
