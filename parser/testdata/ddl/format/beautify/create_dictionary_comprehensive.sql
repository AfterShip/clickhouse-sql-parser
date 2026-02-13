-- Origin SQL:
CREATE OR REPLACE DICTIONARY test.comprehensive_dict 
UUID '12345678-1234-1234-1234-123456789012'
ON CLUSTER production_cluster
(
    id UInt64,
    name String DEFAULT '',
    value Float64 EXPRESSION toFloat64OrZero(name),
    parent_id UInt64 HIERARCHICAL,
    is_active UInt8 INJECTIVE,
    object_id UInt64 IS_OBJECT_ID
)
PRIMARY KEY id
SOURCE(MYSQL(
    host 'localhost'
    port 3306
    user 'root'
    password 'secret'
    db 'test_db'
    table 'dictionary_table'
))
LIFETIME(MIN 1000 MAX 2000)
LAYOUT(HASHED())
SETTINGS(max_block_size = 8192, max_insert_block_size = 1048576);

-- Beautify SQL:
CREATE OR REPLACE DICTIONARY test.comprehensive_dict UUID '12345678-1234-1234-1234-123456789012' ON CLUSTER production_cluster (
  id UInt64,
  name String DEFAULT '',
  value Float64 EXPRESSION toFloat64OrZero(name),
  parent_id UInt64 HIERARCHICAL,
  is_active UInt8 INJECTIVE,
  object_id UInt64 IS_OBJECT_ID
)
PRIMARY KEY id
SOURCE(MYSQL(host 'localhost' port 3306 user 'root' password 'secret' db 'test_db' table 'dictionary_table'))
LIFETIME(MIN 1000 MAX 2000)
LAYOUT(HASHED())
SETTINGS(max_block_size=8192, max_insert_block_size=1048576);
