# Table: openstack_flavors

This table shows data for Openstack Flavors.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_flavors:
  - [openstack_flavor_accesses](openstack_flavor_accesses.md)
  - [openstack_flavor_extra_specs](openstack_flavor_extra_specs.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|id|String|
|disk|Int|
|ram|Int|
|name|String|
|rxtx_factor|Float|
|vcpus|Int|
|is_public|Bool|
|ephemeral|Int|
|description|String|