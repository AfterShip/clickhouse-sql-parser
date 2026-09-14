-- Origin SQL:
SELECT -tuple(1, 2).1, -arr[1], -x::Int64, - -x;
SELECT a = b IN (1), a LIKE b IN (1), a BETWEEN b = c AND d = e;
SELECT a::Int64[1], a::Array(Int64)[1], a::Tuple(Int64, Int64).1;


-- Beautify SQL:
SELECT
  -tuple(1, 2).1,
  -arr[1],
  -x::Int64,
  - -x;
SELECT
  a = b IN (1),
  a LIKE b IN (1),
  a BETWEEN b = c AND d = e;
SELECT
  a::Int64[1],
  a::Array(Int64)[1],
  a::Tuple(Int64, Int64).1;
