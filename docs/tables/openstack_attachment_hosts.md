# Table: openstack_attachment_hosts

This table shows data for Openstack Attachment Hosts.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_attachments](openstack_attachments.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|host|String|
|port|Int|