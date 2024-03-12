# Table: openstack_blockstorage_quotasets_usage

This table shows data for Openstack Blockstorage Quotasets Usage.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|volumes_in_use|`int64`|
|volumes_allocated|`int64`|
|volumes_reserved|`int64`|
|volumes_limit|`int64`|
|snapshots_in_use|`int64`|
|snapshots_allocated|`int64`|
|snapshots_reserved|`int64`|
|snapshots_limit|`int64`|
|gigabytes_in_use|`int64`|
|gigabytes_allocated|`int64`|
|gigabytes_reserved|`int64`|
|gigabytes_limit|`int64`|
|per_volume_gigabytes_in_use|`int64`|
|per_volume_gigabytes_allocated|`int64`|
|per_volume_gigabytes_reserved|`int64`|
|per_volume_gigabytes_limit|`int64`|
|backups_in_use|`int64`|
|backups_allocated|`int64`|
|backups_reserved|`int64`|
|backups_limit|`int64`|
|backup_gigabytes_in_use|`int64`|
|backup_gigabytes_allocated|`int64`|
|backup_gigabytes_reserved|`int64`|
|backup_gigabytes_limit|`int64`|
|groups_in_use|`int64`|
|groups_allocated|`int64`|
|groups_reserved|`int64`|
|groups_limit|`int64`|
|id|`utf8`|