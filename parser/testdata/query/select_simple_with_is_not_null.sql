SELECT f0,f1,f2,f3 as a0
FROM test.events_local
WHERE (f0 IN ('foo', 'bar', 'test'))
  AND (f1 = 'testing')
  AND f2 IS NULL
  AND f3 IS NOT NULL