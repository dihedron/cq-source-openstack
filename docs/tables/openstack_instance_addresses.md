# Table: openstack_instance_addresses

This table shows data for Openstack Instance Addresses.

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
|network|String|
|mac_address|String|
|type|String|
|ip_address|String|
|ip_version|Int|