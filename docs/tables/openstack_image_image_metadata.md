# Table: openstack_image_image_metadata

This table shows data for Openstack Image Image Metadata.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_image_images](openstack_image_images.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|key|`utf8`|
|value|`utf8`|