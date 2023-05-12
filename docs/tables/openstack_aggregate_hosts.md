# Table: openstack_aggregate_hosts

This table shows data for Openstack Aggregate Hosts.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_aggregates](openstack_aggregates.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|name|String|