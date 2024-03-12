# Table: openstack_identity_domain_groups

This table shows data for Openstack Identity Domain Groups.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_identity_domains](openstack_identity_domains.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|description|`utf8`|
|domain_id|`utf8`|
|id|`utf8`|
|name|`utf8`|