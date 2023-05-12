# Table: openstack_instance_flavors

This table shows data for Openstack Instance Flavors.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_instances](openstack_instances.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|name|String|
|vcpus|Int|
|vgpus|Int|
|cores|Int|
|sockets|Int|
|ram|Int|
|disk|Int|
|swap|Int|
|ephemeral|Int|
|rng_allowed|Bool|
|watchdog_action|String|