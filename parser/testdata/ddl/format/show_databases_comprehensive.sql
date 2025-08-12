-- Origin SQL:
SHOW DATABASES LIKE 'prod%' LIMIT 5 INTO OUTFILE '/tmp/prod_dbs.txt' FORMAT JSON

-- Format SQL:
SHOW DATABASES LIKE 'prod%' LIMIT 5 INTO OUTFILE '/tmp/prod_dbs.txt' FORMAT 'JSON';
