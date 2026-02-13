-- Origin SQL:
CREATE OR REPLACE VIEW asdf AS
SELECT id,
       price * 1.5 AS computed_value,
       row_number() OVER (
           PARTITION BY category
           ORDER BY created_at
           RANGE BETWEEN 3600 PRECEDING AND CURRENT ROW
           )       AS rn
FROM source_table
WHERE date >= '2023-01-01';


-- Beautify SQL:
CREATE OR REPLACE VIEW asdf AS SELECT
  id,
  price * 1.5 AS computed_value,
  row_number() OVER (PARTITION BY category ORDER BY
    created_at RANGE BETWEEN 3600 PRECEDING AND CURRENT ROW) AS rn
FROM
  source_table
WHERE
  date >= '2023-01-01';
