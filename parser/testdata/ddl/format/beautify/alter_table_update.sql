-- Origin SQL:
ALTER TABLE test.users UPDATE status = 'active', updated_at = now() WHERE status = 'pending';


-- Beautify SQL:
ALTER TABLE test.users
UPDATE status = 'active', updated_at = now() WHERE status = 'pending';
