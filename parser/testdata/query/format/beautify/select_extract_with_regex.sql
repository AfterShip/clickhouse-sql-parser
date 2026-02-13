-- Origin SQL:
SELECT
  COUNT(1), SRC_TYPE, NODE_CLASS, PORT, CLIENT_PORT
FROM
  test.table
WHERE
  app_id = 999118646
  AND toUnixTimestamp(timestamp) >= 1740366695
  AND toUnixTimestamp(timestamp) <= 1740377495
GROUP BY
  CASE
    WHEN length(extract(instance, '((\\d+\\.){3}\\d+)')) > 0 THEN instance
    ELSE '空'
  END,
  CASE
    WHEN length(extract(client_ip, '((\\d+\\.){3}\\d+)')) > 0 THEN client_ip
    ELSE '空'
  END,
  src_type,
  node_class,
  port,
  client_port
LIMIT 10000


-- Beautify SQL:
SELECT
  COUNT(1),
  SRC_TYPE,
  NODE_CLASS,
  PORT,
  CLIENT_PORT
FROM
  test.table
WHERE
  app_id = 999118646
AND
  toUnixTimestamp(timestamp) >= 1740366695
AND
  toUnixTimestamp(timestamp) <= 1740377495
GROUP BY
  CASE
    WHEN length(EXTRACT(instance, '((\\d+\\.){3}\\d+)')) > 0 THEN instance
    ELSE '空'
  END, CASE
    WHEN length(EXTRACT(client_ip, '((\\d+\\.){3}\\d+)')) > 0 THEN client_ip
    ELSE '空'
  END, src_type, node_class, port, client_port
LIMIT 10000;
