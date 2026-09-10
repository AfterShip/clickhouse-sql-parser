#!/usr/bin/clickhouse
SELECT `a``b`, `a\`b`, "a""b", "a\"b";
SELECT $$hello$$, $tag$it's\n$tag$, $1_$a$$b$1_$, $$$$;
SELECT /* outer /* inner */ outer */ 1 # comment
+ 2;
