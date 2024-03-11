# Table: openstack_compute_aggregate_hosts

This table shows data for Openstack Compute Aggregate Hosts.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_aggregates](openstack_compute_aggregates.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|name|`utf8`|