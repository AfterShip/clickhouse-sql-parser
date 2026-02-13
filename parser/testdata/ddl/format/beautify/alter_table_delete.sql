-- Origin SQL:
ALTER TABLE test.events DELETE WHERE created_at < '2023-01-01';


-- Beautify SQL:
ALTER TABLE test.events
DELETE WHERE created_at < '2023-01-01';
