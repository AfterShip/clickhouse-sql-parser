ALTER TABLE test.users UPDATE status = 'active' IN PARTITION ID '202401' WHERE id > 100;
