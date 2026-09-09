-- Origin SQL:
CREATE TABLE t
(
    `value_1` Dynamic(max_types = 16),
    `value_2` Dynamic(max_types = 16)
)
ENGINE = Memory;


-- Format SQL:
CREATE TABLE t (`value_1` Dynamic(max_types=16), `value_2` Dynamic(max_types=16)) ENGINE = Memory;
