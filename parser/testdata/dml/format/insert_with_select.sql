-- Origin SQL:
INSERT INTO test.visits_null
SELECT
    CounterID,
    StartDate,
    Sign,
    UserID
FROM test.visits;

-- Format SQL:
INSERT INTO test.visits_null SELECT CounterID, StartDate, Sign, UserID FROM test.visits;
