ALTER TABLE test.users UPDATE status = 'active', updated_at = now() WHERE status = 'pending';
