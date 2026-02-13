-- Origin SQL:
ALTER TABLE test.users UPDATE status = 'inactive' IN PARTITION '2024-01-01' WHERE status = 'active';


-- Beautify SQL:
ALTER TABLE test.users
UPDATE status = 'inactive' IN PARTITION '2024-01-01' WHERE status = 'active';
