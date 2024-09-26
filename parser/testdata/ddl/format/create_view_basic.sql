-- Origin SQL:
CREATE VIEW IF NOT EXISTS my_view(col1 String, col2 String)
AS
SELECT
    id,
    name
FROM
    my_table;

-- Format SQL:
CREATE VIEW IF NOT EXISTS my_view (col1 String, col2 String) AS SELECT id, name FROM my_table;
