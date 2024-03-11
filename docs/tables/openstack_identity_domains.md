# Table: openstack_identity_domains

This table shows data for Openstack Identity Domains.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_identity_domains:
  - [openstack_identity_domain_groups](openstack_identity_domain_groups.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|description|`utf8`|
|enabled|`bool`|
|id|`utf8`|
|name|`utf8`|