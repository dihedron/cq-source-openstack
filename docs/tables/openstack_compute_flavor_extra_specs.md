# Table: openstack_compute_flavor_extra_specs

This table shows data for Openstack Compute Flavor Extra Specs.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_flavors](openstack_compute_flavors.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|key|`utf8`|
|value|`utf8`|