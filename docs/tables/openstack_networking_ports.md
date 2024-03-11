# Table: openstack_networking_ports

This table shows data for Openstack Networking Ports.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|ip_addresses|`list<item: utf8, nullable>`|
|ip_address|`utf8`|
|id (PK)|`utf8`|
|network_id|`utf8`|
|name|`utf8`|
|description|`utf8`|
|admin_state_up|`bool`|
|status|`utf8`|
|mac_address|`utf8`|
|fixed_ips|`json`|
|tenant_id|`utf8`|
|project_id|`utf8`|
|device_owner|`utf8`|
|security_groups|`list<item: utf8, nullable>`|
|device_id|`utf8`|
|allowed_address_pairs|`json`|
|tags|`list<item: utf8, nullable>`|
|propagate_uplink_status|`bool`|
|value_specs|`json`|
|revision_number|`int64`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|