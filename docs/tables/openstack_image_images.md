# Table: openstack_image_images

This table shows data for Openstack Image Images.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_image_images:
  - [openstack_image_image_members](openstack_image_image_members.md)
  - [openstack_image_image_metadata](openstack_image_image_metadata.md)
  - [openstack_image_image_properties](openstack_image_image_properties.md)
  - [openstack_image_image_tags](openstack_image_image_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|id|`utf8`|
|name|`utf8`|
|status|`utf8`|
|tags|`list<item: utf8, nullable>`|
|container_format|`utf8`|
|disk_format|`utf8`|
|min_disk|`int64`|
|min_ram|`int64`|
|owner|`utf8`|
|protected|`bool`|
|visibility|`utf8`|
|os_hidden|`bool`|
|checksum|`utf8`|
|metadata|`json`|
|properties|`json`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|
|file|`utf8`|
|schema|`utf8`|
|virtual_size|`int64`|