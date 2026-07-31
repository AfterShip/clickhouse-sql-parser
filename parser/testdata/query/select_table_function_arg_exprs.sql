SELECT * FROM numbers(1 + 1);
SELECT * FROM numbers(intDiv(a, b) + 1);
SELECT * FROM numbers(greatest(1, a - b));
SELECT * FROM cluster('c', numbers(1 + 1));
SELECT * FROM numbers(x AND y);
SELECT * FROM numbers(1.5);
SELECT * FROM numbers(now());
