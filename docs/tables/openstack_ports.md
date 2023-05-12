# Table: openstack_ports

This table shows data for Openstack Ports.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|ip_addresses|StringArray|
|ip_address|String|
|id (PK)|String|
|network_id|String|
|name|String|
|description|String|
|admin_state_up|Bool|
|status|String|
|mac_address|String|
|fixed_ips|JSON|
|tenant_id|String|
|project_id|String|
|device_owner|String|
|security_groups|StringArray|
|device_id|String|
|allowed_address_pairs|JSON|
|tags|StringArray|
|propagate_uplink_status|Bool|
|value_specs|JSON|
|revision_number|Int|
|created_at|Timestamp|
|updated_at|Timestamp|