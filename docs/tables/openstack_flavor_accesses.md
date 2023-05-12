# Table: openstack_flavor_accesses

This table shows data for Openstack Flavor Accesses.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_flavors](openstack_flavors.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|flavor_id|String|
|project_id|String|