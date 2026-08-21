CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id SET x = sum(x);

CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id;

CREATE TABLE t (id UInt64, created DateTime, x UInt64, total UInt64) ENGINE = MergeTree() ORDER BY (id, created) TTL created + INTERVAL 1 DAY GROUP BY id, created SET x = sum(x), total = count();

CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id SET x = sum(x), created + INTERVAL 2 DAY DELETE;

CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY ALL SET x = sum(x);

CREATE TABLE t (id UInt64, created DateTime, x UInt64, y UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id SET x = max(y) > 0;