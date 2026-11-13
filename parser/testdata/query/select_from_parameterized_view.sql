SELECT id
FROM db.v(min_id = {min_id:UInt64}, ids = [{id:UInt64}], since = toDateTime64({since:String}, 3)) AS a
INNER JOIN `db.v`(min_id = 1, ids = [1, 2], since = now()) AS b ON a.id = b.id
