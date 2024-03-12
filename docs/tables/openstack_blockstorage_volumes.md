# Table: openstack_blockstorage_volumes

This table shows data for Openstack Blockstorage Volumes.

The primary key for this table is **id**.

## Relations

The following tables depend on openstack_blockstorage_volumes:
  - [openstack_blockstorage_volumes_backups](openstack_blockstorage_volumes_backups.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|id (PK)|`utf8`|
|status|`utf8`|
|size|`int64`|
|availability_zone|`utf8`|
|created_at|`json`|
|updated_at|`json`|
|attachments|`json`|
|name|`utf8`|
|description|`utf8`|
|volume_type|`utf8`|
|snapshot_id|`utf8`|
|source_volid|`utf8`|
|backup_id|`utf8`|
|group_id|`utf8`|
|metadata|`json`|
|user_id|`utf8`|
|bootable|`utf8`|
|encrypted|`bool`|
|replication_status|`utf8`|
|consistencygroup_id|`utf8`|
|multiattach|`bool`|
|volume_image_metadata|`json`|
|migration_status|`utf8`|
|host|`utf8`|
|migration_status_name|`utf8`|
|migration_status_tenant|`utf8`|
|provider_id|`utf8`|
|service_uuid|`utf8`|
|shared_targets|`bool`|