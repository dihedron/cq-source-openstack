# Table: openstack_blockstorage_quotasets

This table shows data for Openstack Blockstorage Quotasets.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|id|`utf8`|
|volumes|`int64`|
|snapshots|`int64`|
|gigabytes|`int64`|
|per_volume_gigabytes|`int64`|
|backups|`int64`|
|backup_gigabytes|`int64`|
|groups|`int64`|