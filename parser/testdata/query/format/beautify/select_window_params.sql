-- Origin SQL:
-- Parameters in WHERE and in window frames (UInt32 & String; both spacing styles; shorthand frame)
SELECT sum(x) OVER (ORDER BY y ROWS BETWEEN {start:UInt32} PRECEDING AND CURRENT ROW)                           AS total1,
       avg(x) OVER (ORDER BY y ROWS BETWEEN CURRENT ROW AND {end:UInt32} FOLLOWING)                             AS avg1,
       count(*) OVER (ORDER BY y RANGE BETWEEN {range_start:UInt32} PRECEDING AND {range_end:UInt32} FOLLOWING) AS cnt1,
       sum(x) OVER (ROWS {window_size :UInt32} PRECEDING)                                                       AS rows_shorthand
FROM t
WHERE category = {category :String}
  AND type = {type:String};


-- Beautify SQL:
SELECT
  sum(x) OVER (ORDER BY
    y ROWS BETWEEN {start: UInt32} PRECEDING AND CURRENT ROW) AS total1,
  avg(x) OVER (ORDER BY
    y ROWS BETWEEN CURRENT ROW AND {end: UInt32} FOLLOWING) AS avg1,
  count(*) OVER (ORDER BY
    y RANGE BETWEEN {range_start: UInt32} PRECEDING AND {range_end: UInt32} FOLLOWING) AS cnt1,
  sum(x) OVER (ROWS {window_size: UInt32} PRECEDING) AS rows_shorthand
FROM
  t
WHERE
  category = {category: String}
AND
  type = {type: String};
