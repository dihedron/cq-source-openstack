# Table: openstack_compute_aggregates

This table shows data for Openstack Compute Aggregates.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_compute_aggregates:
  - [openstack_compute_aggregate_hosts](openstack_compute_aggregate_hosts.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|availability_zone|`utf8`|
|hosts|`list<item: utf8, nullable>`|
|id|`int64`|
|metadata|`json`|
|name|`utf8`|
|deleted|`bool`|