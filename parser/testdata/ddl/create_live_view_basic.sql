CREATE LIVE VIEW my_live_view
WITH TIMEOUT 10 TO my_destination(id String)
AS SELECT id FROM my_table;
