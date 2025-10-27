ALTER TABLE test.users UPDATE status = 'active' IN PARTITION '2024-01' WHERE id > 100;
