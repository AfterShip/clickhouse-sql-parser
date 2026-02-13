-- Origin SQL:
CREATE DICTIONARY test.my_dict (
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
    user 'default'
    password ''
    db 'test'
    table 'dict_table'
))
LIFETIME(MIN 1000 MAX 2000)
LAYOUT(HASHED())
SETTINGS(max_block_size = 8192);

-- Beautify SQL:
CREATE DICTIONARY test.my_dict (
  id UInt64,
  name String DEFAULT '',
  value Float64 EXPRESSION toFloat64OrZero(name),
  parent_id UInt64 HIERARCHICAL,
  is_active UInt8 INJECTIVE,
  object_id UInt64 IS_OBJECT_ID
)
PRIMARY KEY id
SOURCE(MYSQL(host 'localhost' port 3306 user 'default' password '' db 'test' table 'dict_table'))
LIFETIME(MIN 1000 MAX 2000)
LAYOUT(HASHED())
SETTINGS(max_block_size=8192);
