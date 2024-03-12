# Table: openstack_blockstorage_limits

This table shows data for Openstack Blockstorage Limits.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|max_total_volumes|`int64`|
|max_total_snapshots|`int64`|
|max_total_volume_gigabytes|`int64`|
|max_total_backups|`int64`|
|max_total_backup_gigabytes|`int64`|
|total_volumes_used|`int64`|
|total_gigabytes_used|`int64`|
|total_snapshots_used|`int64`|
|total_backups_used|`int64`|
|total_backup_gigabytes_used|`int64`|
|regex|`list<item: utf8, nullable>`|
|uri|`list<item: utf8, nullable>`|
|verb|`list<item: utf8, nullable>`|
|next_available|`list<item: utf8, nullable>`|
|unit|`list<item: utf8, nullable>`|
|value|`list<item: int64, nullable>`|
|remaining|`list<item: int64, nullable>`|