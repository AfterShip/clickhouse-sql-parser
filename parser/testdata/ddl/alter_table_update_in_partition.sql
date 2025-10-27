ALTER TABLE test.users UPDATE status = 'active' WHERE id > 100 IN PARTITION '2024-01';
