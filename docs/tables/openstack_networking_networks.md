# Table: openstack_networking_networks

This table shows data for Openstack Networking Networks.

The primary key for this table is **id**.

## Relations

The following tables depend on openstack_networking_networks:
  - [openstack_networking_network_subnets](openstack_networking_network_subnets.md)
  - [openstack_networking_network_tags](openstack_networking_network_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|id (PK)|`utf8`|
|name|`utf8`|
|description|`utf8`|
|admin_state_up|`bool`|
|status|`utf8`|
|subnets|`list<item: utf8, nullable>`|
|tenant_id|`utf8`|
|project_id|`utf8`|
|shared|`bool`|
|availability_zone_hints|`list<item: utf8, nullable>`|
|tags|`list<item: utf8, nullable>`|
|revision_number|`int64`|