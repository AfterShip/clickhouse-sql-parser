-- Origin SQL:
INSERT INTO t0(user_id, message, timestamp, metric) VALUES
    (?, ?, ?, ?),
    (?, ?, ?, ?),
    (?, ?, ?, ?),
    (?, ?, ?, ?)
;

INSERT INTO test_with_typed_columns (id, created_at)
VALUES ({id: Int32}, {created_at: DateTime64(6)});

-- Format SQL:
INSERT INTO t0 (user_id, message, timestamp, metric) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?);
INSERT INTO test_with_typed_columns (id, created_at) VALUES ({id:Int32}, {created_at:DateTime64(6)});
