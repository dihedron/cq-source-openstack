# Table: openstack_networking_security_groups

This table shows data for Openstack Networking Security Groups.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|security_group_rule_ids|`list<item: utf8, nullable>`|
|id (PK)|`utf8`|
|name|`utf8`|
|description|`utf8`|
|security_group_rules|`json`|
|tenant_id|`utf8`|
|project_id|`utf8`|
|tags|`list<item: utf8, nullable>`|