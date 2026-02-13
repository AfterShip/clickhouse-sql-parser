-- Origin SQL:
WITH
    monthly AS (
        SELECT toStartOfMonth(date) AS month,
               department,
               avg(salary)          AS avg_salary
        FROM salary_table
        WHERE year = 2023
        GROUP BY month, department
    ),
    ranked AS (
        SELECT month,
               department,
               avg_salary,
               row_number() OVER (PARTITION BY department ORDER BY avg_salary DESC) AS dept_rank
        FROM monthly
    )
SELECT month,
       department,
       avg_salary,
       lag(avg_salary, 1, 0) OVER (
           PARTITION BY department
           ORDER BY month
           ROWS BETWEEN 1 PRECEDING AND CURRENT ROW
           ) AS prev_month_avg
FROM ranked
WHERE dept_rank <= 5
ORDER BY month, department;


-- Beautify SQL:
WITH
  monthly AS (SELECT
    toStartOfMonth(date) AS month,
    department,
    avg(salary) AS avg_salary
  FROM
    salary_table
  WHERE
    year = 2023
  GROUP BY
    month, department),
  ranked AS (SELECT
    month,
    department,
    avg_salary,
    row_number() OVER (PARTITION BY department ORDER BY
      avg_salary DESC) AS dept_rank
  FROM
    monthly)
SELECT
  month,
  department,
  avg_salary,
  lag(avg_salary, 1, 0) OVER (PARTITION BY department ORDER BY
    month ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) AS prev_month_avg
FROM
  ranked
WHERE
  dept_rank <= 5
ORDER BY
  month,
  department;
