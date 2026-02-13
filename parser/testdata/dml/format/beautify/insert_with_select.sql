-- Origin SQL:
INSERT INTO test.visits_null
SELECT
    CounterID,
    StartDate,
    Sign,
    UserID
FROM test.visits;

-- Beautify SQL:
INSERT INTO test.visits_null
SELECT
  CounterID,
  StartDate,
  Sign,
  UserID
FROM
  test.visits;
