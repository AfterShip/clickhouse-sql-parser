CREATE VIEW IF NOT EXISTS cluster_name.my_view
        UUID '3493e374-e2bb-481b-b493-e374e2bb981b'
        ON CLUSTER 'my_cluster'
AS (
    SELECT
    column1,
    column2
    FROM
    my_other_table
);