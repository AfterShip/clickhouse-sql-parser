-- Origin SQL:
CREATE NAMED COLLECTION IF NOT EXISTS my_collection ON CLUSTER my_cluster
AS key1 = 'value1' OVERRIDABLE,
key2 = 'value2' NOT OVERRIDABLE,
key3 = 'value3';


-- Format SQL:
CREATE NAMED COLLECTION IF NOT EXISTS my_collection ON CLUSTER my_cluster AS key1 = 'value1' OVERRIDABLE, key2 = 'value2' NOT OVERRIDABLE, key3 = 'value3';
