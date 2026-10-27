-- Origin SQL:
SELECT ip.address FROM public_ips ip WHERE ip.region = 'eastus';
SELECT * FROM t1 ip FINAL LEFT JOIN t2 key ON ip.a = key.a;
SELECT * FROM numbers(3) ip;


-- Beautify SQL:
SELECT
  ip.address
FROM
  public_ips AS ip
WHERE
  ip.region = 'eastus';
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
