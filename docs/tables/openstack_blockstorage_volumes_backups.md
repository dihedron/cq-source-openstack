# Table: openstack_blockstorage_volumes_backups

This table shows data for Openstack Blockstorage Volumes Backups.

The primary key for this table is **id**.

## Relations

This table depends on [openstack_blockstorage_volumes](openstack_blockstorage_volumes.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|id (PK)|`utf8`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|
|name|`utf8`|
|description|`utf8`|
|volume_id|`utf8`|
|snapshot_id|`utf8`|
|status|`utf8`|
|size|`int64`|
|object_count|`int64`|
|container|`utf8`|
|has_dependent_backups|`bool`|
|fail_reason|`utf8`|
|is_incremental|`bool`|
|data_timestamp|`timestamp[us, tz=UTC]`|
|project_id|`utf8`|
|metadata|`json`|
|availability_zone|`utf8`|