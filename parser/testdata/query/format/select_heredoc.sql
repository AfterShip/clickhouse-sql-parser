-- Origin SQL:
SELECT $tag$hello$tag$;
SELECT $x$ AS a FROM t WHERE b = $x$;
SELECT $$it's\raw$$ AS service$$name, 1 AS a$b, 2 AS $open$;


-- Format SQL:
SELECT 'hello';
SELECT ' AS a FROM t WHERE b = ';
SELECT 'it''s\\raw' AS service$$name, 1 AS a$b, 2 AS $open$;
