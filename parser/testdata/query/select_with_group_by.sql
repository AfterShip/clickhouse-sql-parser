SELECT
    datacenter,
    distro,
    SUM (quantity) AS qty
FROM
    servers
GROUP BY
    GROUPING SETS(
    (datacenter,distro),
    (datacenter),
    (distro),
    ()
);

SELECT
    datacenter,
    distro,
    SUM (quantity) AS qty
FROM
    servers
GROUP BY ALL;