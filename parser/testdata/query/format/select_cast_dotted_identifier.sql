-- Origin SQL:
-- Test CAST with various dotted identifier lengths
SELECT CAST(column, 'String') AS single_part;
SELECT CAST(table.column, 'String') AS two_parts;
SELECT CAST(db.table.column, 'String') AS three_parts;
SELECT CAST(some.long.json.path, 'String') AS four_parts;
SELECT CAST(a.very.long.nested.json.path.with.many.parts, 'String') AS many_parts;

-- Format SQL:
SELECT CAST(column, 'String') AS single_part;
SELECT CAST(table.column, 'String') AS two_parts;
SELECT CAST(db.table.column, 'String') AS three_parts;
SELECT CAST(some.long.json.path, 'String') AS four_parts;
SELECT CAST(a.very.long.nested.json.path.with.many.parts, 'String') AS many_parts;
