# Table: openstack_compute_instance_flavors

This table shows data for Openstack Compute Instance Flavors.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_instances](openstack_compute_instances.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|name|`utf8`|
|vcpus|`int64`|
|vgpus|`int64`|
|cores|`int64`|
|sockets|`int64`|
|ram|`int64`|
|disk|`int64`|
|swap|`int64`|
|ephemeral|`int64`|
|rng_allowed|`bool`|
|watchdog_action|`utf8`|