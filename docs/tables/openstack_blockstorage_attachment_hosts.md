# Table: openstack_blockstorage_attachment_hosts

This table shows data for Openstack Blockstorage Attachment Hosts.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_blockstorage_attachments](openstack_blockstorage_attachments.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|host|`utf8`|
|port|`int64`|