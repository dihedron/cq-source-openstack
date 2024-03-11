# Table: openstack_compute_instance_flavor_extra_specs

This table shows data for Openstack Compute Instance Flavor Extra Specs.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_instances](openstack_compute_instances.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|key|`utf8`|
|value|`utf8`|