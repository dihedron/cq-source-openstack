# Table: openstack_flavor_extra_specs

This table shows data for Openstack Flavor Extra Specs.

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
|key|String|
|value|String|