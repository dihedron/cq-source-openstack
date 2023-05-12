# Table: openstack_attachments

This table shows data for Openstack Attachments.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_attachments:
  - [openstack_attachment_hosts](openstack_attachment_hosts.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|attached_at|Timestamp|
|detached_at|Timestamp|
|access_mode|String|
|attach_mode|String|
|attachment_id|String|
|auth_enabled|Bool|
|auth_username|String|
|cluster_name|String|
|discard|Bool|
|driver_volume_type|String|
|encrypted|Bool|
|hosts|StringArray|
|keyring|Bool|
|name|String|
|ports|StringArray|
|secret_type|String|
|secret_uuid|String|
|volume_id|String|
|id|String|
|instance_id|String|
|status|String|
|project_id|String|
|connection_info|JSON|