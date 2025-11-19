-- CREATE TABLE with columns AS table function (remoteSecure)
CREATE TABLE test_remote
(
    id UInt64,
    name String,
    value Int32
)
AS remoteSecure('host.example.com', 'source_db', 'source_table', 'user', 'password');

-- Simpler test case with remote()
CREATE TABLE test_table (id UInt64, name String) AS remote('localhost', 'db', 'source_table');
