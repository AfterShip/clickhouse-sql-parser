-- Origin SQL:
SELECT $tag$hello$tag$;
SELECT $x$ AS a FROM t WHERE b = $x$;
SELECT $$it's\raw$$ AS service$$name, 1 AS a$b, 2 AS $open$;
SELECT $t_1$it's\raw
and $$ inside$t_1$;


-- Beautify SQL:
SELECT
  $tag$hello$tag$;
SELECT
  $x$ AS a FROM t WHERE b = $x$;
SELECT
  $$it's\raw$$ AS service$$name,
  1 AS a$b,
  2 AS $open$;
SELECT
  $t_1$it's\raw
and $$ inside$t_1$;
