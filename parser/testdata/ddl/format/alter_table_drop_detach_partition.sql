-- Origin SQL:
ALTER TABLE app_utc_00.app_message_as_notification_organization_sent_stats_i_d_local DROP DETACHED PARTITION '2022-05-24' SETTINGS allow_drop_detached = 1;

-- Format SQL:
ALTER TABLE app_utc_00.app_message_as_notification_organization_sent_stats_i_d_local
DROP DETACHED PARTITION '2022-05-24' SETTINGS allow_drop_detached=1;
