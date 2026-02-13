-- Origin SQL:
-- Comprehensive window spec coverage: ad-hoc specs, named windows, ROWS/RANGE frames
-- multiple functions, multi-column PARTITION/ORDER, and expression-based specs.
SELECT
    -- Ad-hoc windows
    sum(x) OVER (ORDER BY y ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)                                         AS running_total,
    avg(x)
        OVER (PARTITION BY z ORDER BY y RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)                                    AS avg_range1,

    -- Named window reuse (OVER w1)
    row_number() OVER w1                                                                                              AS rn_w1,
    rank() OVER w1                                                                                                    AS rank_w1,
    sum(x) OVER w1                                                                                                    AS sum_w1,

    -- Frame variants (incl. shorthand & RANGE)
    sum(x) OVER (ROWS 10 PRECEDING)                                                                                   AS rows_10_preceding,
    sum(x) OVER (ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)                                                    AS rows_cur_to_unbounded_following,
    sum(x) OVER (ROWS BETWEEN 5 PRECEDING AND 3 FOLLOWING)                                                            AS rows_5p_3f,
    sum(x) OVER (RANGE BETWEEN 10 PRECEDING AND CURRENT ROW)                                                          AS range_10p_cur,
    sum(x) OVER (RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)                                                   AS range_unbounded_to_cur,

    -- Ranking & navigation
    row_number() OVER (PARTITION BY y ORDER BY x)                                                                     AS row_num,
    rank() OVER (PARTITION BY y ORDER BY x)                                                                           AS rank_val,
    dense_rank() OVER (PARTITION BY y ORDER BY x)                                                                     AS dense_rank_val,
    first_value(x)
                OVER (PARTITION BY y ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)             AS first_val,
    last_value(x)
               OVER (PARTITION BY y ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)              AS last_val,
    lag(x, 1) OVER (PARTITION BY y ORDER BY x)                                                                        AS prev_x,
    lead(x, 1) OVER (PARTITION BY y ORDER BY x)                                                                       AS next_x,
    percent_rank() OVER (PARTITION BY y ORDER BY x)                                                                   AS pct_rank,

    -- Named window reference via OVER w (no parentheses)
    sum(x) OVER w                                                                                                     AS sum_over_w,
    avg(x) OVER w                                                                                                     AS avg_over_w,
    row_number() OVER w                                                                                               AS rn_over_w,

    -- Multiple columns in PARTITION BY / ORDER BY
    count(*) OVER (PARTITION BY col1, col2, col3 ORDER BY col4, col5 DESC)                                            AS cnt_multi,
    sum(val)
        OVER (PARTITION BY col1, col2 ORDER BY col4, col5 ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)           AS total_multi,

    -- Expressions in PARTITION/ORDER
    sum(amount)
        OVER (PARTITION BY date_trunc('day', timestamp) ORDER BY timestamp ROWS BETWEEN 10 PRECEDING AND CURRENT ROW) AS daily_total,
    avg(amount)
        OVER (ORDER BY extract(HOUR FROM timestamp) RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)                        AS hourly_avg
FROM t
WINDOW w AS (ORDER BY y),
       w1 AS (PARTITION BY y ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW),
       w4 AS (PARTITION BY y ORDER BY x ROWS BETWEEN 3 PRECEDING AND CURRENT ROW),
       w5 AS (PARTITION BY z ORDER BY x RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW);


-- Beautify SQL:
SELECT
  sum(x) OVER (ORDER BY
    y ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
  avg(x) OVER (PARTITION BY z ORDER BY
    y RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS avg_range1,
  row_number() OVER w1 AS rn_w1,
  rank() OVER w1 AS rank_w1,
  sum(x) OVER w1 AS sum_w1,
  sum(x) OVER (ROWS 10 PRECEDING) AS rows_10_preceding,
  sum(x) OVER (ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) AS rows_cur_to_unbounded_following,
  sum(x) OVER (ROWS BETWEEN 5 PRECEDING AND 3 FOLLOWING) AS rows_5p_3f,
  sum(x) OVER (RANGE BETWEEN 10 PRECEDING AND CURRENT ROW) AS range_10p_cur,
  sum(x) OVER (RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS range_unbounded_to_cur,
  row_number() OVER (PARTITION BY y ORDER BY
    x) AS row_num,
  rank() OVER (PARTITION BY y ORDER BY
    x) AS rank_val,
  dense_rank() OVER (PARTITION BY y ORDER BY
    x) AS dense_rank_val,
  first_value(x) OVER (PARTITION BY y ORDER BY
    x ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) AS first_val,
  last_value(x) OVER (PARTITION BY y ORDER BY
    x ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) AS last_val,
  lag(x, 1) OVER (PARTITION BY y ORDER BY
    x) AS prev_x,
  lead(x, 1) OVER (PARTITION BY y ORDER BY
    x) AS next_x,
  percent_rank() OVER (PARTITION BY y ORDER BY
    x) AS pct_rank,
  sum(x) OVER w AS sum_over_w,
  avg(x) OVER w AS avg_over_w,
  row_number() OVER w AS rn_over_w,
  count(*) OVER (PARTITION BY col1, col2, col3 ORDER BY
    col4,
    col5 DESC) AS cnt_multi,
  sum(val) OVER (PARTITION BY col1, col2 ORDER BY
    col4,
    col5 ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS total_multi,
  sum(amount) OVER (PARTITION BY date_trunc('day', timestamp) ORDER BY
    timestamp ROWS BETWEEN 10 PRECEDING AND CURRENT ROW) AS daily_total,
  avg(amount) OVER (ORDER BY
    EXTRACT(HOUR FROM timestamp) RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS hourly_avg
FROM
  t
WINDOW w AS (ORDER BY
  y), w1 AS (PARTITION BY y ORDER BY
  x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW), w4 AS (PARTITION BY y ORDER BY
  x ROWS BETWEEN 3 PRECEDING AND CURRENT ROW), w5 AS (PARTITION BY z ORDER BY
  x RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW);
