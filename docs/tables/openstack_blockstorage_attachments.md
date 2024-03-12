# Table: openstack_blockstorage_attachments

This table shows data for Openstack Blockstorage Attachments.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_blockstorage_attachments:
  - [openstack_blockstorage_attachment_hosts](openstack_blockstorage_attachment_hosts.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|attached_at|`timestamp[us, tz=UTC]`|
|detached_at|`timestamp[us, tz=UTC]`|
|access_mode|`utf8`|
|attach_mode|`utf8`|
|attachment_id|`utf8`|
|auth_enabled|`bool`|
|auth_username|`utf8`|
|cluster_name|`utf8`|
|discard|`bool`|
|driver_volume_type|`utf8`|
|encrypted|`bool`|
|hosts|`list<item: utf8, nullable>`|
|keyring|`bool`|
|name|`utf8`|
|ports|`list<item: utf8, nullable>`|
|secret_type|`utf8`|
|secret_uuid|`utf8`|
|volume_id|`utf8`|
|id|`utf8`|
|instance_id|`utf8`|
|status|`utf8`|
|project_id|`utf8`|
|connection_info|`json`|