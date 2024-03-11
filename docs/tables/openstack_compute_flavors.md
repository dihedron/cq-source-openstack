# Table: openstack_compute_flavors

This table shows data for Openstack Compute Flavors.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_compute_flavors:
  - [openstack_compute_flavor_accesses](openstack_compute_flavor_accesses.md)
  - [openstack_compute_flavor_extra_specs](openstack_compute_flavor_extra_specs.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|id|`utf8`|
|disk|`int64`|
|ram|`int64`|
|name|`utf8`|
|rxtx_factor|`float64`|
|vcpus|`int64`|
|is_public|`bool`|
|ephemeral|`int64`|
|description|`utf8`|