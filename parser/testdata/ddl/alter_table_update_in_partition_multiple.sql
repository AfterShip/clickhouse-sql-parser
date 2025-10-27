ALTER TABLE test.users UPDATE status = 'active', updated_at = now() WHERE id > 100 IN PARTITION ('2024-01');
