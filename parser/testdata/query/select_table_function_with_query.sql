SELECT 1, (SELECT 70) AS `power`, number
FROM
numbers(
    plus(
        ifNull((SELECT 1 AS bin_count, 1),
        1)
    )
)