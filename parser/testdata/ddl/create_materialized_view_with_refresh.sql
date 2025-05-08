CREATE MATERIALIZED VIEW fresh_mv
REFRESH EVERY 1 HOUR OFFSET 10 MINUTE
RANDOMIZE FOR 1 SECOND
DEPENDS ON  table_v5
SETTINGS
    randomize_for = 1,
    randomize_offset = 10,
    randomize_period = 1
APPEND TO target_table_name
EMPTY
AS SELECT
    `field_1`,
    `field_2`,
    `field_3`,
FROM table_v5