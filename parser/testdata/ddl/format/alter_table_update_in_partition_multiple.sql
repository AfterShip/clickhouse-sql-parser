-- Origin SQL:
ALTER TABLE test.users UPDATE status = 'active', updated_at = now() IN PARTITION ('2024-01') WHERE id > 1000;


-- Format SQL:
ALTER TABLE test.users UPDATE status = 'active', updated_at = now() IN PARTITION ('2024-01') WHERE id > 1000;
