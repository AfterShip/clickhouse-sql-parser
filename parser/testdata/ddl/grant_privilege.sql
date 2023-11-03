GRANT SELECT(x,y) ON db.table TO john;
GRANT SELECT(x,y) ON db.table TO john WITH GRANT OPTION WITH ADMIN OPTION;
GRANT SELECT(x,y) ON db.* TO john;
GRANT SELECT(x,y) ON *.table TO john;
GRANT SELECT(x,y) ON *.* TO john;
GRANT SELECT(x,y) ON *.table TO CURRENT_USER;
GRANT SELECT(x,y) ON *.table TO CURRENT_USER,john,mary;
GRANT ALL ON *.* TO admin_role WITH GRANT OPTION;
