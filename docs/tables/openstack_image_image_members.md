# Table: openstack_image_image_members

This table shows data for Openstack Image Image Members.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_image_images](openstack_image_images.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|created_at|`timestamp[us, tz=UTC]`|
|image_id|`utf8`|
|member_id|`utf8`|
|schema|`utf8`|
|status|`utf8`|
|updated_at|`timestamp[us, tz=UTC]`|