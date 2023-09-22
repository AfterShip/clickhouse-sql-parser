create table test.event_all
ON CLUSTER 'default_cluster'
AS test.evnets_local
ENGINE = Distributed(
    default_cluster,
    test,
    events_local,
    rand()
) SETTINGS fsync_after_insert=0;
