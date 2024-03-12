# Table: openstack_compute_instance_addresses

This table shows data for Openstack Compute Instance Addresses.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_instances](openstack_compute_instances.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|network|`utf8`|
|mac_address|`utf8`|
|type|`utf8`|
|ip_address|`utf8`|
|ip_version|`int64`|