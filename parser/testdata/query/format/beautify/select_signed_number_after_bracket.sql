-- Origin SQL:
SELECT (1)-1;
SELECT (1)+1.5;
SELECT arr[1]-1;
SELECT (toUnixTimestamp(now())-3600)*1000000000;


-- Beautify SQL:
SELECT
  (1) - 1;
SELECT
  (1) + 1.5;
SELECT
  arr[1] - 1;
SELECT
  (toUnixTimestamp(now()) - 3600) * 1000000000;
