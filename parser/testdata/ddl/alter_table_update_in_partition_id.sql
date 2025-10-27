ALTER TABLE test.users UPDATE status = 'active' WHERE id > 100 IN PARTITION ID '202401';
