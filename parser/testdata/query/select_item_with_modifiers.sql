SELECT c0 REPLACE(c0 AS c1) FROM t0;
SELECT * REPLACE(i + 1 AS i) FROM t1;
SELECT * REPLACE(i + 1 AS i) EXCEPT (j) APPLY(sum) from t2;