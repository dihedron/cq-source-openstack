# Table: openstack_baremetal_ports

This table shows data for Openstack Baremetal Ports.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|id|`utf8`|
|uuid|`utf8`|
|tenant_id|`utf8`|
|project_id|`utf8`|
|device_id|`utf8`|
|network_id|`utf8`|
|ip_addresses|`list<item: utf8, nullable>`|
|ip_address|`utf8`|
|name|`utf8`|
|description|`utf8`|
|admin_state_up|`bool`|
|status|`utf8`|
|mac_address|`utf8`|
|fixed_ips|`list<item: utf8, nullable>`|
|device_owner|`utf8`|
|security_groups|`json`|
|allowed_address_pairs|`list<item: utf8, nullable>`|
|tags|`list<item: utf8, nullable>`|
|propagate_uplink_status|`bool`|
|value_specs|`utf8`|
|revision_number|`int64`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|