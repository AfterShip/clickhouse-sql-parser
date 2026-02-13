-- Origin SQL:
ALTER TABLE example_table MODIFY SETTING max_part_loading_threads=8, max_parts_in_total=50000;

-- Beautify SQL:
ALTER TABLE example_table
MODIFY SETTING max_part_loading_threads=8, max_parts_in_total=50000;
