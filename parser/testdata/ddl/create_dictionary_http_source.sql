CREATE DICTIONARY asns (asn UInt32, name String)
PRIMARY KEY asn
SOURCE(HTTP(URL 'http://o/asns.csv' FORMAT 'CSVWithNames'))
LIFETIME(MIN 0 MAX 3600)
LAYOUT(HASHED());

CREATE DICTIONARY os (id UInt32, name String)
PRIMARY KEY id
SOURCE(HTTP(
    url 'http://[::1]/os.tsv'
    format 'TabSeparated'
    credentials(user 'user' password 'password')
    headers(header(name 'API-KEY' value 'key'))
))
LIFETIME(MIN 0 MAX 3600)
LAYOUT(HASHED());
