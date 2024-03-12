# Table: openstack_networking_security_group_rules

This table shows data for Openstack Networking Security Group Rules.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|id (PK)|`utf8`|
|direction|`utf8`|
|description|`utf8`|
|ethertype|`utf8`|
|security_group_id|`utf8`|
|port_range_min|`int64`|
|port_range_max|`int64`|
|protocol|`utf8`|
|remote_group_id|`utf8`|
|remote_ip_prefix|`utf8`|
|tenant_id|`utf8`|
|project_id|`utf8`|