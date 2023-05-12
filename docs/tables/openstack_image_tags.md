# Table: openstack_image_tags

This table shows data for Openstack Image Tags.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_images](openstack_images.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|value|String|