-- Origin SQL:
ALTER TABLE test.users UPDATE status = 'inactive' IN PARTITION '2024-01-01' WHERE status = 'active';


-- Format SQL:
ALTER TABLE test.users UPDATE status = 'inactive' IN PARTITION '2024-01-01' WHERE status = 'active';
