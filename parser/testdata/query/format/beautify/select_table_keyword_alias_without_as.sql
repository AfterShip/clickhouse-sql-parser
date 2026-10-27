-- Origin SQL:
SELECT ip.name FROM azure_network_public_ip_addresses ip WHERE ip.location = 'eastus';
SELECT * FROM t1 ip FINAL LEFT JOIN t2 key ON ip.a = key.a;
SELECT * FROM numbers(3) ip;


-- Beautify SQL:
SELECT
  ip.name
FROM
  azure_network_public_ip_addresses AS ip
WHERE
  ip.location = 'eastus';
SELECT
  *
FROM
  t1 AS ip FINAL
  LEFT JOIN
    t2 AS key ON ip.a = key.a;
SELECT
  *
FROM
  numbers(3) AS ip;
