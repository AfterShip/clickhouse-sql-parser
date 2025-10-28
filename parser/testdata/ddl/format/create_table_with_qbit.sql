-- Origin SQL:
CREATE TABLE test.qbit_example (
    id UInt32,
    vec QBit(Float32, 8)
) ENGINE = Memory;


-- Format SQL:
CREATE TABLE test.qbit_example (id UInt32, vec QBit(Float32, 8)) ENGINE = Memory;
