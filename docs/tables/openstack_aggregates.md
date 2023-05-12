# Table: openstack_aggregates

This table shows data for Openstack Aggregates.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_aggregates:
  - [openstack_aggregate_hosts](openstack_aggregate_hosts.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|availability_zone|String|
|hosts|StringArray|
|id|Int|
|metadata|JSON|
|name|String|
|deleted|Bool|