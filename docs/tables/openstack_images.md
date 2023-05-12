# Table: openstack_images

This table shows data for Openstack Images.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_images:
  - [openstack_image_metadata](openstack_image_metadata.md)
  - [openstack_image_properties](openstack_image_properties.md)
  - [openstack_image_tags](openstack_image_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|id|String|
|name|String|
|status|String|
|tags|StringArray|
|container_format|String|
|disk_format|String|
|min_disk|Int|
|min_ram|Int|
|owner|String|
|protected|Bool|
|visibility|String|
|os_hidden|Bool|
|checksum|String|
|metadata|JSON|
|properties|JSON|
|created_at|Timestamp|
|updated_at|Timestamp|
|file|String|
|schema|String|
|virtual_size|Int|