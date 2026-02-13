-- Origin SQL:
SELECT DISTINCT ON(album,artist) record_id FROM records

-- Beautify SQL:
SELECT DISTINCT ON (album, artist)
  record_id
FROM
  records;
