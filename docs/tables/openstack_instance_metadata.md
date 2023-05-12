# Table: openstack_instance_metadata

This table shows data for Openstack Instance Metadata.

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
|key|String|
|value|String|