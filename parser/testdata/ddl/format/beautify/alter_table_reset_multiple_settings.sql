-- Origin SQL:
ALTER TABLE example_table RESET SETTING max_part_loading_threads, max_parts_in_total, another_setting;

-- Beautify SQL:
ALTER TABLE example_table
RESET SETTING max_part_loading_threads, max_parts_in_total, another_setting;
