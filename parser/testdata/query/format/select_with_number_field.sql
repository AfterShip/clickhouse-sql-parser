-- Origin SQL:
SELECT foo, bar.1, foo.2 FROM foo ARRAY JOIN m as bar

-- Format SQL:
SELECT foo, bar.1, foo.2 FROM foo  ARRAY JOIN m AS bar;
