# Table: openstack_compute_flavor_accesses

This table shows data for Openstack Compute Flavor Accesses.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_compute_flavors](openstack_compute_flavors.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|flavor_id|`utf8`|
|project_id|`utf8`|