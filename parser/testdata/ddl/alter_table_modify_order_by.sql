ALTER TABLE flows MODIFY ORDER BY (TimeReceived, ExporterAddress);
ALTER TABLE flows MODIFY ORDER BY TimeReceived;
ALTER TABLE db.flows ON CLUSTER c MODIFY ORDER BY (TimeReceived, toStartOfHour(TimeReceived));
ALTER TABLE flows ADD COLUMN ExporterAddress String, MODIFY ORDER BY (TimeReceived, ExporterAddress);
