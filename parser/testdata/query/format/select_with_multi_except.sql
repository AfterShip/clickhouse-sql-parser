-- Origin SQL:
SELECT number FROM numbers(1, 10) EXCEPT SELECT number FROM numbers(3, 6) EXCEPT SELECT number FROM numbers(8, 9)

-- Format SQL:
SELECT number FROM numbers(1, 10) EXCEPT SELECT number FROM numbers(3, 6) EXCEPT SELECT number FROM numbers(8, 9);
