-- Origin SQL:
CREATE VIEW IF NOT EXISTS db.my_view
(
    `id` Int64,
    `name` String
)
COMMENT '{"blueprint_hash":"abc123"}'
AS SELECT
    id,
    name
FROM db.my_table;


-- Format SQL:
CREATE VIEW IF NOT EXISTS db.my_view (`id` Int64, `name` String) COMMENT '{"blueprint_hash":"abc123"}' AS SELECT id, name FROM db.my_table;
