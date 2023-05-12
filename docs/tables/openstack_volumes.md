# Table: openstack_volumes

This table shows data for Openstack Volumes.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|id (PK)|String|
|status|String|
|size|Int|
|availability_zone|String|
|created_at|JSON|
|updated_at|JSON|
|attachments|JSON|
|name|String|
|description|String|
|volume_type|String|
|snapshot_id|String|
|source_volid|String|
|backup_id|String|
|group_id|String|
|metadata|JSON|
|user_id|String|
|bootable|String|
|encrypted|Bool|
|replication_status|String|
|consistencygroup_id|String|
|multiattach|Bool|
|volume_image_metadata|JSON|
|migration_status|String|
|host|String|
|migration_status_name|String|
|migration_status_tenant|String|
|provider_id|String|
|service_uuid|String|
|shared_targets|Bool|